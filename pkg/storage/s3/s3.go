package s3

import (
	"context"
)

type S3Storage struct {
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
}

func (s *S3Storage) Read(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}

func NewS3Storage(region, bucket, accessKeyID, secretAccessKey string) (*S3Storage, error) {
	return &S3Storage{
		Region:          region,
		Bucket:          bucket,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}, nil
}
