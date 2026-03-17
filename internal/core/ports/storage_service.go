package ports

import (
	"context"
	"mime/multipart"
)

type StorageService interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader) (string, error)
}
