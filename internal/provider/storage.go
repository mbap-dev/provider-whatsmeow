package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Storage encapsulates an S3-compatible object store used to persist incoming media
// files.  Files are stored in the configured bucket and a presigned URL is
// returned to the caller for temporary public access.
type Storage struct {
	client *minio.Client
	bucket string
	expiry time.Duration
}

// NewStorage creates a new S3/MinIO client using the provided credentials.
// The bucket is created if it does not already exist.  The resulting Storage
// will generate presigned URLs valid for the specified expiry duration.
func NewStorage(endpoint, accessKey, secretKey, bucket, region string, useSSL bool, expiry time.Duration) (*Storage, error) {
	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
		Region: region,
	})
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	exists, err := cli.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := cli.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: region}); err != nil {
			return nil, err
		}
	}
	return &Storage{client: cli, bucket: bucket, expiry: expiry}, nil
}

// Upload stores the given data under objectName and returns a presigned URL for retrieval.
func (s *Storage) Upload(ctx context.Context, objectName string, data []byte, contentType string) (string, error) {
	reader := bytes.NewReader(data)
	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, int64(len(data)), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}
	u, err := s.client.PresignedGetObject(ctx, s.bucket, objectName, s.expiry, url.Values{})
	if err != nil {
		return "", err
	}
	// store a metadata JSON under v13.0/<id> so clients expecting the
	// WhatsApp Cloud API can resolve the downloadable URL by ID.
	id := strings.TrimSuffix(path.Base(objectName), path.Ext(objectName))
	metaObj := path.Join("v13.0", id)
	meta := map[string]string{"url": u.String()}
	buf, err := json.Marshal(meta)
	if err == nil {
		_, err = s.client.PutObject(ctx, s.bucket, metaObj, bytes.NewReader(buf), int64(len(buf)), minio.PutObjectOptions{ContentType: "application/json"})
	}
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
