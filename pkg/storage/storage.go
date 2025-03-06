package storage

import (
	"context"
	"io"
)

// todo: implement the storage interface, we can use local file system, or remote storage like s3/nfs/...
type Storage struct {
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) Read(ctx context.Context, key string) (io.ReadCloser, error) {
	return nil, nil
}

func (s *Storage) Write(ctx context.Context, key string, data []byte) error {
	return nil
}
