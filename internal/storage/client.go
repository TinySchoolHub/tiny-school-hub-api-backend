package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	appconfig "github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/config"
	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/core/domain"
)

const (
	MaxFileSize     = 5 * 1024 * 1024 // 5 MB
	PresignedExpiry = 15 * time.Minute
	DownloadExpiry  = 1 * time.Hour
)

var AllowedContentTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

// Client wraps S3-compatible storage operations
type Client struct {
	s3Client      *s3.Client
	presignClient *s3.PresignClient
	bucket        string
}

// NewClient creates a new storage client using AWS SDK v2
func NewClient(cfg *appconfig.StorageConfig) (*Client, error) {
	// Build endpoint URL
	endpointURL := cfg.Endpoint
	if cfg.Insecure {
		endpointURL = "http://" + cfg.Endpoint
	} else if !strings.HasPrefix(endpointURL, "http") {
		endpointURL = "https://" + cfg.Endpoint
	}

	// Create AWS SDK v2 config
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKey,
			cfg.SecretKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client with custom endpoint
	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpointURL)
		o.UsePathStyle = cfg.UsePathStyle
	})

	presignClient := s3.NewPresignClient(s3Client)

	return &Client{
		s3Client:      s3Client,
		presignClient: presignClient,
		bucket:        cfg.Bucket,
	}, nil
}

// ValidateContentType checks if the content type is allowed
func ValidateContentType(contentType string) error {
	if !AllowedContentTypes[contentType] {
		return domain.ErrInvalidFileType
	}
	return nil
}

// ValidateFileSize checks if the file size is within limits
func ValidateFileSize(size int) error {
	if size > MaxFileSize {
		return domain.ErrFileTooLarge
	}
	return nil
}

// GeneratePresignedPutURL generates a presigned URL for uploading
func (c *Client) GeneratePresignedPutURL(ctx context.Context, key, contentType string) (string, error) {
	if err := ValidateContentType(contentType); err != nil {
		return "", err
	}

	presignResult, err := c.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(PresignedExpiry))

	if err != nil {
		return "", fmt.Errorf("failed to presign PUT request: %w", err)
	}

	return presignResult.URL, nil
}

// GeneratePresignedGetURL generates a presigned URL for downloading
func (c *Client) GeneratePresignedGetURL(ctx context.Context, key string) (string, error) {
	presignResult, err := c.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(DownloadExpiry))

	if err != nil {
		return "", fmt.Errorf("failed to presign GET request: %w", err)
	}

	return presignResult.URL, nil
}

// DeleteObject deletes an object from storage
func (c *Client) DeleteObject(ctx context.Context, key string) error {
	_, err := c.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

// ObjectExists checks if an object exists
func (c *Client) ObjectExists(ctx context.Context, key string) (bool, error) {
	_, err := c.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") || strings.Contains(err.Error(), "NoSuchKey") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// HealthCheck performs a health check on the storage
func (c *Client) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	maxKeys := int32(1)
	_, err := c.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:  aws.String(c.bucket),
		MaxKeys: &maxKeys,
	})
	if err != nil {
		return fmt.Errorf("storage health check failed: %w", err)
	}
	return nil
}
