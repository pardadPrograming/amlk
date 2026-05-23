package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"amlakcrm/backend/internal/domain"
	"amlakcrm/backend/internal/repository"
	"amlakcrm/backend/internal/support"
)

const (
	FileStatusTemporary = "temporary"
	FileStatusAttached  = "attached"

	UploadPurposeChannelMedia  = "channel_media"
	UploadPurposeVaultFile     = "vault_file"
	UploadPurposePropertyMedia = "property_media"
	UploadPurposeBusinessLogo  = "business_logo"

	defaultUploadTTL = time.Hour
)

type UploadService struct {
	store     repository.Store
	objectDir string
}

type UploadInput struct {
	Purpose    string
	TargetType string
	TargetID   string
	BusinessID string
}

type UploadResult struct {
	File domain.FileObject `json:"file"`
}

type uploadPolicy struct {
	imageMaxBytes int64
	videoMaxBytes int64
	fileMaxBytes  int64
	imageOnly     bool
	prefix        string
}

func NewUploadService(store repository.Store, objectDir string) *UploadService {
	return &UploadService{store: store, objectDir: objectDir}
}

func (s *UploadService) Upload(ctx context.Context, user domain.User, input UploadInput, header *multipart.FileHeader) (UploadResult, error) {
	if header == nil || header.Size <= 0 {
		return UploadResult{}, errors.New("فایل معتبر نیست")
	}
	input.Purpose = strings.TrimSpace(input.Purpose)
	input.TargetType = strings.TrimSpace(input.TargetType)
	input.TargetID = strings.TrimSpace(input.TargetID)
	input.BusinessID = strings.TrimSpace(input.BusinessID)
	if err := s.authorizeUpload(ctx, user.ID, input); err != nil {
		return UploadResult{}, err
	}
	policy, err := s.policy(input)
	if err != nil {
		return UploadResult{}, err
	}
	body, kind, contentType, ext, err := s.prepareUpload(ctx, header, policy)
	if err != nil {
		return UploadResult{}, err
	}
	key := fmt.Sprintf("%s/%s%s", strings.Trim(policy.prefix, "/"), support.NewID(), ext)
	if err := s.writeObject(key, body); err != nil {
		return UploadResult{}, err
	}
	file, err := s.store.CreateFile(ctx, domain.FileObject{
		OwnerID:     input.TargetID,
		UploaderID:  user.ID,
		Purpose:     input.Purpose,
		TargetType:  input.TargetType,
		TargetID:    input.TargetID,
		Status:      FileStatusTemporary,
		Provider:    "s3-compatible",
		Bucket:      "amlak",
		Key:         key,
		URL:         "/objects/" + key,
		Kind:        kind,
		ContentType: contentType,
		Size:        int64(len(body)),
		ExpiresAt:   time.Now().UTC().Add(defaultUploadTTL),
	})
	if err != nil {
		_ = os.Remove(filepath.Join(s.objectDir, filepath.FromSlash(key)))
		return UploadResult{}, err
	}
	return UploadResult{File: file}, nil
}

func (s *UploadService) Claim(ctx context.Context, userID string, fileID string, input UploadInput, ttl time.Duration) (domain.FileObject, error) {
	file, err := s.store.GetFile(ctx, fileID)
	if err != nil {
		return domain.FileObject{}, err
	}
	if file.Status != "" && file.Status != FileStatusTemporary && file.Status != FileStatusAttached {
		return domain.FileObject{}, errors.New("فایل در وضعیت قابل استفاده نیست")
	}
	if !file.ExpiresAt.IsZero() && time.Now().UTC().After(file.ExpiresAt) {
		return domain.FileObject{}, errors.New("مهلت استفاده از فایل تمام شده است")
	}
	if file.UploaderID != "" && file.UploaderID != userID {
		return domain.FileObject{}, errors.New("فایل متعلق به کاربر دیگری است")
	}
	if file.Purpose != "" && file.Purpose != input.Purpose {
		return domain.FileObject{}, errors.New("هدف فایل با درخواست همخوانی ندارد")
	}
	if file.TargetID != "" && file.TargetID != input.TargetID {
		return domain.FileObject{}, errors.New("مقصد فایل با درخواست همخوانی ندارد")
	}
	if err := s.authorizeUpload(ctx, userID, input); err != nil {
		return domain.FileObject{}, err
	}
	file.Purpose = input.Purpose
	file.TargetType = input.TargetType
	file.TargetID = input.TargetID
	file.OwnerID = input.TargetID
	file.Status = FileStatusAttached
	if ttl > 0 {
		file.ExpiresAt = time.Now().UTC().Add(ttl)
	}
	return s.store.UpdateFile(ctx, file)
}

func (s *UploadService) StartExpiredFileCleanup(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = time.Hour
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		s.cleanupExpiredFiles(ctx)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.cleanupExpiredFiles(ctx)
			}
		}
	}()
}

func (s *UploadService) cleanupExpiredFiles(ctx context.Context) {
	for {
		files, err := s.store.ListExpiredFiles(ctx, time.Now().UTC(), 100)
		if err != nil || len(files) == 0 {
			return
		}
		for _, file := range files {
			for _, key := range append([]string{file.Key}, file.PreloadKeys...) {
				if strings.TrimSpace(key) == "" {
					continue
				}
				_ = os.Remove(filepath.Join(s.objectDir, filepath.FromSlash(key)))
			}
			_ = s.store.DeleteFile(ctx, file.ID)
		}
		if len(files) < 100 {
			return
		}
	}
}

func (s *UploadService) authorizeUpload(ctx context.Context, userID string, input UploadInput) error {
	switch input.Purpose {
	case UploadPurposeChannelMedia, UploadPurposeVaultFile:
		return s.authorizeChannelWrite(ctx, userID, input.TargetID)
	case UploadPurposePropertyMedia:
		if input.BusinessID == "" || input.TargetID == "" {
			return errors.New("شناسه املاک و فایل ملکی برای آپلود الزامی است")
		}
		member, err := s.store.GetMemberByUser(ctx, input.BusinessID, userID)
		if err != nil || member.Status != domain.MemberActive {
			return errors.New("دسترسی آپلود فایل ملکی وجود ندارد")
		}
		if _, err := s.store.GetPropertyFile(ctx, input.BusinessID, input.TargetID); err != nil {
			return errors.New("فایل ملکی مقصد پیدا نشد")
		}
		return nil
	case UploadPurposeBusinessLogo:
		member, err := s.store.GetMemberByUser(ctx, input.TargetID, userID)
		if err != nil || member.Status != domain.MemberActive {
			return errors.New("دسترسی آپلود لوگوی املاک وجود ندارد")
		}
		if member.Role == domain.RoleOwner || member.Role == domain.RoleManager || domain.HasPermission(member, domain.PermBusinessUpdate) {
			return nil
		}
		return errors.New("دسترسی آپلود لوگوی املاک وجود ندارد")
	default:
		return errors.New("هدف آپلود معتبر نیست")
	}
}

func (s *UploadService) authorizeChannelWrite(ctx context.Context, userID, channelID string) error {
	channel, err := s.store.GetChannel(ctx, channelID)
	if err != nil {
		return errors.New("کانال مقصد پیدا نشد")
	}
	if channel.OwnerUserID == userID {
		return nil
	}
	if channel.BusinessID != "" {
		if member, err := s.store.GetMemberByUser(ctx, channel.BusinessID, userID); err == nil && member.Status == domain.MemberActive {
			if channel.Type != domain.ChannelTypeBusinessVault ||
				member.Role == domain.RoleOwner ||
				member.Role == domain.RoleManager ||
				domain.HasPermission(member, domain.PermBusinessUpdate) ||
				domain.HasPermission(member, domain.PermMembersManage) ||
				s.isChannelAdmin(ctx, userID, channelID) {
				return nil
			}
		}
	}
	if channel.Type == domain.ChannelTypeUserVault && s.isChannelAdmin(ctx, userID, channelID) {
		return nil
	}
	if _, err := s.store.GetChannelMember(ctx, channelID, userID); err == nil && channel.Type != domain.ChannelTypeUserVault && channel.Type != domain.ChannelTypeBusinessVault {
		return nil
	}
	return errors.New("دسترسی آپلود در مقصد وجود ندارد")
}

func (s *UploadService) isChannelAdmin(ctx context.Context, userID, channelID string) bool {
	member, err := s.store.GetChannelMember(ctx, channelID, userID)
	return err == nil && member.Status == domain.ChannelMemberActive && member.Role == domain.ChannelMemberRoleAdmin
}

func (s *UploadService) policy(input UploadInput) (uploadPolicy, error) {
	prefixTarget := input.TargetID
	if prefixTarget == "" {
		prefixTarget = "general"
	}
	switch input.Purpose {
	case UploadPurposeChannelMedia:
		return uploadPolicy{imageMaxBytes: maxChannelImageBytes, videoMaxBytes: maxChannelVideoBytes, fileMaxBytes: maxChannelVideoBytes, prefix: "channels/" + prefixTarget}, nil
	case UploadPurposeVaultFile:
		return uploadPolicy{imageMaxBytes: maxChannelImageBytes, videoMaxBytes: maxChannelVideoBytes, fileMaxBytes: maxChannelVideoBytes, prefix: "vaults/" + prefixTarget}, nil
	case UploadPurposePropertyMedia:
		return uploadPolicy{imageMaxBytes: 500 * 1024, videoMaxBytes: maxChannelVideoBytes, fileMaxBytes: maxChannelVideoBytes, prefix: "properties/" + input.BusinessID + "/" + prefixTarget}, nil
	case UploadPurposeBusinessLogo:
		return uploadPolicy{imageMaxBytes: 500 * 1024, fileMaxBytes: 2 * 1024 * 1024, imageOnly: true, prefix: "business-logos/" + prefixTarget}, nil
	default:
		return uploadPolicy{}, errors.New("هدف آپلود معتبر نیست")
	}
}

func (s *UploadService) prepareUpload(ctx context.Context, header *multipart.FileHeader, policy uploadPolicy) ([]byte, string, string, string, error) {
	contentType := uploadContentType(header)
	switch {
	case strings.HasPrefix(contentType, "image/"):
		src, err := header.Open()
		if err != nil {
			return nil, "", "", "", err
		}
		defer src.Close()
		body, err := optimizeUploadImage(src, policy.imageMaxBytes)
		if err != nil {
			return nil, "", "", "", err
		}
		return body, "image", "image/jpeg", ".jpg", nil
	case strings.HasPrefix(contentType, "video/"):
		if policy.imageOnly {
			return nil, "", "", "", errors.New("برای این مقصد فقط تصویر مجاز است")
		}
		body, err := optimizeUploadVideo(ctx, header, policy.videoMaxBytes)
		if err != nil {
			return nil, "", "", "", err
		}
		return body, "video", "video/mp4", ".mp4", nil
	default:
		if policy.imageOnly {
			return nil, "", "", "", errors.New("برای این مقصد فقط تصویر مجاز است")
		}
		if header.Size > policy.fileMaxBytes {
			return nil, "", "", "", errors.New("حجم فایل بیش از حد مجاز است")
		}
		body, err := readUpload(header)
		if err != nil {
			return nil, "", "", "", err
		}
		ext := filepath.Ext(header.Filename)
		if ext == "" {
			ext = ".bin"
		}
		return body, "file", contentType, ext, nil
	}
}

func optimizeUploadImage(src multipart.File, maxBytes int64) ([]byte, error) {
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, errors.New("تصویر معتبر نیست")
	}
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	current := img
	var encoded []byte
	for scale := 100; scale >= 35; scale -= 10 {
		if scale != 100 {
			current = uploadResizeNearest(img, width*scale/100, height*scale/100)
		}
		for quality := 82; quality >= 38; quality -= 6 {
			var buf bytes.Buffer
			if err := jpeg.Encode(&buf, current, &jpeg.Options{Quality: quality}); err != nil {
				return nil, err
			}
			encoded = buf.Bytes()
			if int64(len(encoded)) <= maxBytes {
				return encoded, nil
			}
		}
	}
	return nil, errors.New("تصویر بعد از فشرده‌سازی هنوز بزرگ‌تر از حد مجاز است")
}

func optimizeUploadVideo(ctx context.Context, header *multipart.FileHeader, maxBytes int64) ([]byte, error) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return nil, errors.New("ffmpeg برای فشرده‌سازی ویدئو لازم است")
	}
	tempDir, err := os.MkdirTemp("", "amlak-upload-video-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)
	inputPath := filepath.Join(tempDir, "input"+filepath.Ext(header.Filename))
	outputPath := filepath.Join(tempDir, "output.mp4")
	if err := saveMultipartUpload(header, inputPath); err != nil {
		return nil, err
	}
	attempts := []struct {
		scale string
		crf   string
	}{
		{scale: "scale=-2:720", crf: "28"},
		{scale: "scale=-2:540", crf: "32"},
		{scale: "scale=-2:480", crf: "36"},
	}
	var lastOutput string
	for _, attempt := range attempts {
		_ = os.Remove(outputPath)
		cmd := exec.CommandContext(ctx, "ffmpeg", "-y", "-i", inputPath, "-vf", attempt.scale, "-c:v", "libx264", "-preset", "veryfast", "-crf", attempt.crf, "-c:a", "aac", "-b:a", "96k", "-movflags", "+faststart", outputPath)
		output, err := cmd.CombinedOutput()
		lastOutput = strings.TrimSpace(string(output))
		if err != nil {
			continue
		}
		body, err := os.ReadFile(outputPath)
		if err != nil {
			return nil, err
		}
		if int64(len(body)) <= maxBytes {
			return body, nil
		}
	}
	if lastOutput == "" {
		lastOutput = "ویدئو بعد از فشرده‌سازی هنوز بزرگ‌تر از حد مجاز است"
	}
	return nil, errors.New(lastOutput)
}

func uploadResizeNearest(src image.Image, width, height int) image.Image {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	srcBounds := src.Bounds()
	for y := 0; y < height; y++ {
		sy := srcBounds.Min.Y + y*srcBounds.Dy()/height
		for x := 0; x < width; x++ {
			sx := srcBounds.Min.X + x*srcBounds.Dx()/width
			dst.Set(x, y, src.At(sx, sy))
		}
	}
	return dst
}

func uploadContentType(header *multipart.FileHeader) string {
	if header == nil {
		return ""
	}
	if contentType := strings.TrimSpace(header.Header.Get("Content-Type")); contentType != "" && contentType != "application/octet-stream" {
		return contentType
	}
	if ext := filepath.Ext(header.Filename); ext != "" {
		if contentType := mime.TypeByExtension(ext); contentType != "" {
			return contentType
		}
	}
	return header.Header.Get("Content-Type")
}

func saveMultipartUpload(header *multipart.FileHeader, path string) error {
	src, err := header.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}

func (s *UploadService) writeObject(key string, body []byte) error {
	path := filepath.Join(s.objectDir, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, body, 0644)
}
