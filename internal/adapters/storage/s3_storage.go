package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Storage struct {
	client     *minio.Client
	bucketName string
	endpoint   string
}

func NewS3Storage(endpoint, accessKeyID, secretAccessKey, bucketName string) (*S3Storage, error) {
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint URL: %w", err)
	}

	useSSL := parsedURL.Scheme == "https"
	host := parsedURL.Host

	client, err := minio.New(host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
		Region: "auto",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client: %w", err)
	}

	return &S3Storage{
		client:     client,
		bucketName: bucketName,
		endpoint:   endpoint,
	}, nil
}

func (s *S3Storage) UploadImage(ctx context.Context, reader io.Reader, filename string, contentType string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		ext = ".jpg"
	}

	objectName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err := s.client.PutObject(ctx, s.bucketName, objectName, reader, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object: %w", err)
	}

	// Construct public URL
	publicURL := fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucketName, objectName)
	return publicURL, nil
}
