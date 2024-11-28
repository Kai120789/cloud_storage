package minio

import (
	"bytes"
	"cloud/internal/dto"
	"cloud/internal/models"
	"context"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client *minio.Client
	bucket string
}

// Init MinIO client
func NewMinioStorage(endpoint, accessKey, secretKey, bucket string) (*MinioStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	// Create bucket
	err = client.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := client.BucketExists(context.Background(), bucket)
		if errBucketExists == nil && exists {

		} else {
			return nil, err
		}
	}

	minio := MinioStorage{
		client: client,
		bucket: bucket,
	}

	return &minio, nil
}

func (s *MinioStorage) CreateNewFileOrFold(obj dto.Object) (*models.Object, error) {
	_, err := s.client.PutObject(context.Background(), s.bucket, obj.Path, bytes.NewReader(obj.Content), int64(len(obj.Content)), minio.PutObjectOptions{})
	if err != nil {
		return nil, err
	}

	createdObj := &models.Object{
		Name:      obj.Name,
		Path:      obj.Path,
		UserID:    obj.UserID,
		CreatedAt: time.Now(),
	}

	return createdObj, nil
}

func (s *MinioStorage) DeleteItem(path string) error {
	ctx := context.Background()
	err := s.client.RemoveObject(ctx, s.bucket, path, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (s *MinioStorage) RenameItem(obj dto.Object) (*models.Object, error) {
	ctx := context.Background()
	newPath := obj.Path // Новый путь из DTO

	// Копируем объект
	src := minio.CopySrcOptions{Bucket: s.bucket, Object: obj.Path}
	dst := minio.CopyDestOptions{Bucket: s.bucket, Object: newPath}
	_, err := s.client.CopyObject(ctx, dst, src)
	if err != nil {
		return nil, err
	}

	// Удаляем старый объект
	err = s.DeleteItem(obj.Path)
	if err != nil {
		return nil, err
	}

	retObj := models.Object{
		Name:      obj.Name,
		Path:      newPath,
		UserID:    obj.UserID,
		CreatedAt: time.Now(),
	}

	return &retObj, nil
}

func (s *MinioStorage) ListDirectory(path string) ([]models.Object, error) {
	ctx := context.Background()

	objects := []models.Object{}
	objectCh := s.client.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{
		Prefix:    path,
		Recursive: false,
	})

	for obj := range objectCh {
		if obj.Err != nil {
			return nil, obj.Err
		}

		objects = append(objects, models.Object{
			Name:      obj.Key,
			Path:      obj.Key,
			CreatedAt: obj.LastModified,
		})
	}

	return objects, nil
}

func (s *MinioStorage) SearchFiles(query string) ([]models.Object, error) {
	ctx := context.Background()

	objects := []models.Object{}
	objectCh := s.client.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{
		Recursive: true,
	})

	for obj := range objectCh {
		if obj.Err != nil {
			return nil, obj.Err
		}

		if query == "" || contains(obj.Key, query) {
			objects = append(objects, models.Object{
				Name:      obj.Key,
				Path:      obj.Key,
				CreatedAt: obj.LastModified,
			})
		}
	}

	return objects, nil
}

func contains(key, query string) bool {
	return strings.Contains(key, query)
}
