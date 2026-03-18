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
	client         *minio.Client
	bucketName     string
	publicEndpoint string
}

func NewS3Storage(connectionEndpoint, accessKeyID, secretAccessKey, bucketName, publicEndpoint string) (*S3Storage, error) {
	parsedURL, err := url.Parse(connectionEndpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid connection endpoint URL: %w", err)
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

	storage := &S3Storage{
		client:         client,
		bucketName:     bucketName,
		publicEndpoint: publicEndpoint,
	}

	// Ensure bucket exists and is public
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`, bucketName)
	err = client.SetBucketPolicy(ctx, bucketName, policy)
	if err != nil {
		return nil, fmt.Errorf("failed to set bucket policy: %w", err)
	}

	return storage, nil
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

	// Ensure public endpoint doesn't end with a slash for consistent URL construction
	publicEndpoint := strings.TrimSuffix(s.publicEndpoint, "/")

	// Construct public URL
	publicURL := fmt.Sprintf("%s/%s/%s", publicEndpoint, s.bucketName, objectName)
	return publicURL, nil
}
