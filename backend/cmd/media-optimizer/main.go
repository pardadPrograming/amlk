package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"log/slog"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"amlakcrm/backend/internal/config"
)

const (
	defaultMaxImageBytes = 800 * 1024
	defaultMaxVideoBytes = 50 * 1024 * 1024
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("POST /optimize", optimizeHandler)

	server := &http.Server{
		Addr:              cfg.MediaOptimizerAddr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}
	logger.Info("media optimizer listening", "addr", cfg.MediaOptimizerAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("media optimizer failed", "error", err)
		os.Exit(1)
	}
}

func optimizeHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(220 << 20); err != nil {
		writeOptimizerError(w, http.StatusBadRequest, "invalid multipart body")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeOptimizerError(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	contentType := r.FormValue("contentType")
	if contentType == "" {
		contentType = header.Header.Get("Content-Type")
	}
	if (contentType == "" || contentType == "application/octet-stream") && filepath.Ext(header.Filename) != "" {
		if detected := mime.TypeByExtension(filepath.Ext(header.Filename)); detected != "" {
			contentType = detected
		}
	}
	maxImageBytes := formInt64(r, "maxImageBytes", defaultMaxImageBytes)
	maxVideoBytes := formInt64(r, "maxVideoBytes", defaultMaxVideoBytes)

	switch {
	case strings.HasPrefix(contentType, "image/"):
		body, err := optimizeImage(file, maxImageBytes)
		if err != nil {
			writeOptimizerError(w, http.StatusBadRequest, err.Error())
			return
		}
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("X-Media-Kind", "image")
		w.Header().Set("X-Media-Extension", ".jpg")
		_, _ = w.Write(body)
	case strings.HasPrefix(contentType, "video/"):
		body, err := optimizeVideo(r.Context(), file, header, maxVideoBytes)
		if err != nil {
			writeOptimizerError(w, http.StatusBadRequest, err.Error())
			return
		}
		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("X-Media-Kind", "video")
		w.Header().Set("X-Media-Extension", ".mp4")
		_, _ = w.Write(body)
	default:
		writeOptimizerError(w, http.StatusBadRequest, "only image and video media can be optimized")
	}
}

func optimizeImage(src multipart.File, maxBytes int64) ([]byte, error) {
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, errors.New("invalid image")
	}
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	current := img
	var encoded []byte
	for scale := 100; scale >= 35; scale -= 10 {
		if scale != 100 {
			current = resizeNearest(img, width*scale/100, height*scale/100)
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
	return nil, errors.New("image is still larger than the configured maximum after compression")
}

func optimizeVideo(ctx context.Context, src multipart.File, header *multipart.FileHeader, maxBytes int64) ([]byte, error) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return nil, errors.New("ffmpeg is required")
	}
	tempDir, err := os.MkdirTemp("", "amlak-media-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	inputPath := filepath.Join(tempDir, "input"+filepath.Ext(header.Filename))
	outputPath := filepath.Join(tempDir, "output.mp4")
	input, err := os.Create(inputPath)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(input, src); err != nil {
		_ = input.Close()
		return nil, err
	}
	if err := input.Close(); err != nil {
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
		cmd := exec.CommandContext(
			ctx,
			"ffmpeg", "-y", "-i", inputPath,
			"-vf", attempt.scale,
			"-c:v", "libx264", "-preset", "veryfast", "-crf", attempt.crf,
			"-c:a", "aac", "-b:a", "96k",
			"-movflags", "+faststart",
			outputPath,
		)
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
		lastOutput = "video is still larger than the configured maximum after compression"
	}
	return nil, errors.New(lastOutput)
}

func resizeNearest(src image.Image, width, height int) image.Image {
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

func formInt64(r *http.Request, key string, fallback int64) int64 {
	value := strings.TrimSpace(r.FormValue(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func writeOptimizerError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
