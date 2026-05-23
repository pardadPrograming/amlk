package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type OptimizedMedia struct {
	Body        []byte
	ContentType string
	Kind        string
	Extension   string
}

type MediaOptimizerClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewMediaOptimizerClient(baseURL string) *MediaOptimizerClient {
	return &MediaOptimizerClient{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		httpClient: &http.Client{
			Timeout: 2 * time.Minute,
		},
	}
}

func (c *MediaOptimizerClient) Enabled() bool {
	return c != nil && c.baseURL != ""
}

func (c *MediaOptimizerClient) Optimize(ctx context.Context, header *multipart.FileHeader, maxImageBytes, maxVideoBytes int64) (OptimizedMedia, error) {
	if !c.Enabled() {
		return OptimizedMedia{}, errors.New("media optimizer is not configured")
	}
	src, err := header.Open()
	if err != nil {
		return OptimizedMedia{}, err
	}
	defer src.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		return OptimizedMedia{}, err
	}
	if _, err := io.Copy(part, src); err != nil {
		return OptimizedMedia{}, err
	}
	_ = writer.WriteField("maxImageBytes", strconv.FormatInt(maxImageBytes, 10))
	_ = writer.WriteField("maxVideoBytes", strconv.FormatInt(maxVideoBytes, 10))
	if contentType := detectOptimizerUploadContentType(header); contentType != "" {
		_ = writer.WriteField("contentType", contentType)
	}
	if err := writer.Close(); err != nil {
		return OptimizedMedia{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/optimize", &body)
	if err != nil {
		return OptimizedMedia{}, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := c.httpClient.Do(req)
	if err != nil {
		return OptimizedMedia{}, err
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return OptimizedMedia{}, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return OptimizedMedia{}, fmt.Errorf("media optimizer failed: %s", strings.TrimSpace(string(resBody)))
	}
	contentType := res.Header.Get("Content-Type")
	kind := res.Header.Get("X-Media-Kind")
	extension := res.Header.Get("X-Media-Extension")
	if extension == "" {
		extension = filepath.Ext(header.Filename)
	}
	if extension == "" {
		extension = extensionForContentType(contentType)
	}
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}
	return OptimizedMedia{
		Body:        resBody,
		ContentType: contentType,
		Kind:        kind,
		Extension:   extension,
	}, nil
}

func detectOptimizerUploadContentType(header *multipart.FileHeader) string {
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

func extensionForContentType(contentType string) string {
	switch {
	case strings.HasPrefix(contentType, "image/"):
		return ".jpg"
	case strings.HasPrefix(contentType, "video/"):
		return ".mp4"
	case strings.HasPrefix(contentType, "audio/"):
		return ".bin"
	default:
		return ".bin"
	}
}
