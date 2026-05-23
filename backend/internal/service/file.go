package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"

	"amlakcrm/backend/internal/domain"
	"amlakcrm/backend/internal/repository"
	"amlakcrm/backend/internal/support"
)

type FileService struct {
	store repository.Store
}

func NewFileService(store repository.Store) *FileService {
	return &FileService{store: store}
}

func (s *FileService) SaveBusinessLogo(ctx context.Context, ownerID string, header *multipart.FileHeader) (domain.FileObject, error) {
	if header.Size > 2*1024*1024 {
		return domain.FileObject{}, errors.New("حجم لوگو نباید بیشتر از ۲ مگابایت باشد")
	}
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return domain.FileObject{}, errors.New("فقط فایل تصویری مجاز است")
	}
	key := fmt.Sprintf("business-logos/%s-%s", support.NewID(), header.Filename)
	return s.store.CreateFile(ctx, domain.FileObject{
		OwnerID:     ownerID,
		Provider:    "s3-compatible",
		Bucket:      "amlak",
		Key:         key,
		URL:         "/objects/" + key,
		ContentType: contentType,
		Size:        header.Size,
	})
}
