package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"community-server/internal/config"
)

type UploadService struct {
	uploadDir string
	publicURL string
}

func NewUploadService(cfg *config.Config) *UploadService {
	uploadDir := cfg.File.Path
	if uploadDir == "" {
		uploadDir = "./uploads/"
	}
	publicURL := cfg.File.ExternalPath
	if publicURL == "" {
		publicURL = "/uploads/"
	}
	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("创建上传目录失败: %v", err))
	}
	return &UploadService{
		uploadDir: uploadDir,
		publicURL: strings.TrimRight(publicURL, "/"),
	}
}

type UploadResult struct {
	URL      string `json:"url"`
	FileName string `json:"file_name"`
	Size     int64  `json:"size"`
}

func (s *UploadService) Upload(file *multipart.FileHeader) (*UploadResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	subDir := time.Now().Format("2006/01/02")
	dir := filepath.Join(s.uploadDir, subDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("创建子目录失败: %w", err)
	}

	dstPath := filepath.Join(dir, name)
	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	written, err := io.Copy(dst, src)
	if err != nil {
		return nil, fmt.Errorf("写入文件失败: %w", err)
	}

	url := fmt.Sprintf("%s/%s/%s", s.publicURL, subDir, name)
	return &UploadResult{
		URL:      url,
		FileName: file.Filename,
		Size:     written,
	}, nil
}
