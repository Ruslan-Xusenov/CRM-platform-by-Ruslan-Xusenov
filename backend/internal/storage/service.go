package storage

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Service struct {
	client *minio.Client
}

func NewService(endpoint, accessKey, secretKey string, useSSL bool) (*Service, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	// Ensure default buckets exist
	ctx := context.Background()
	for _, bucket := range []string{"call-recordings", "crm-files"} {
		exists, err := client.BucketExists(ctx, bucket)
		if err != nil {
			slog.Warn("Failed to check bucket", "bucket", bucket, "error", err)
			continue
		}
		if !exists {
			if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
				slog.Warn("Failed to create bucket", "bucket", bucket, "error", err)
			} else {
				slog.Info("Created bucket", "bucket", bucket)
			}
		}
	}

	return &Service{client: client}, nil
}

// GetPresignedUploadURL generates a presigned URL for uploading a file.
func (s *Service) GetPresignedUploadURL(bucket, objectName string, expiry time.Duration) (string, error) {
	presignedURL, err := s.client.PresignedPutObject(context.Background(), bucket, objectName, expiry)
	if err != nil {
		return "", fmt.Errorf("presigned upload URL: %w", err)
	}
	return presignedURL.String(), nil
}

// GetPresignedDownloadURL generates a presigned URL for downloading a file.
func (s *Service) GetPresignedDownloadURL(bucket, objectName string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := s.client.PresignedGetObject(context.Background(), bucket, objectName, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("presigned download URL: %w", err)
	}
	return presignedURL.String(), nil
}

// DeleteObject removes a file from storage.
func (s *Service) DeleteObject(bucket, objectName string) error {
	return s.client.RemoveObject(context.Background(), bucket, objectName, minio.RemoveObjectOptions{})
}
