package s3

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"domashka-backend/config"
)

type Client struct {
	config *config.S3Config
	mc     *minio.Client
}

func New(config *config.S3Config) (*Client, error) {
	if config == nil {
		return nil, nil
	}
	mc, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: true,
		Region: config.Region,
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		config: config,
		mc:     mc,
	}, nil
}

func (c *Client) UploadPicture(ctx context.Context, filePrefix string, fileHeader *multipart.FileHeader) (string, error) {
	// Открываем файл
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Генерируем ключ
	ext := filepath.Ext(fileHeader.Filename)
	key := fmt.Sprintf("%s_%d%s", filePrefix, time.Now().Unix(), ext)

	// Загружаем объект
	info, err := c.mc.PutObject(
		ctx,
		c.config.Bucket,
		key,
		src,
		fileHeader.Size,
		minio.PutObjectOptions{
			ContentType: fileHeader.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		return "", fmt.Errorf("minio PutObject: %w", err)
	}
	fmt.Println(info)
	publicURL := fmt.Sprintf("https://%s/%s/%s", c.config.Endpoint, c.config.Bucket, key)
	return publicURL, nil
}
