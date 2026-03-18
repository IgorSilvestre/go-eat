package ports

import (
	"context"
	"io"
)

type StorageService interface {
	UploadImage(ctx context.Context, reader io.Reader, filename string, contentType string) (string, error)
}
