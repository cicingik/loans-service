package storage

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// ConfigStorage ...
type ConfigStorage struct {
	URL       string `envconfig:"url"`
	AccessKey string `envconfig:"access_key"`
	AccessID  string `envconfig:"access_id"`
	UseSSL    bool   `envconfig:"use_ssl"`
	Bucket    string `envconfig:"bucket"`
}

// Storage ...
type Storage struct {
	cfg    *ConfigStorage
	client *minio.Client
}

// NewStorage ...
func NewStorage(cfg *ConfigStorage) (*Storage, error) {
	client, err := minio.New(cfg.URL, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessID, cfg.AccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	return &Storage{
		cfg:    cfg,
		client: client,
	}, nil
}

// Upload ...
func (s *Storage) Upload(ctx context.Context, filePath string) (string, error) {
	objectName := strings.Split(filePath, "/")

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	mimeType := http.DetectContentType(buffer)

	_, err = s.client.FPutObject(
		ctx, s.cfg.Bucket,
		objectName[len(objectName)-1],
		filePath,
		minio.PutObjectOptions{ContentType: mimeType},
	)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
