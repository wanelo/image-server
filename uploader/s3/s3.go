package s3

import (
	"os"
	"path/filepath"
	"time"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
)

// Uploader for S3
type Uploader struct {
	BaseDir string
}

var bucket *s3.Bucket

// Upload copies a file int a bucket in S3
func (u *Uploader) Upload(source string, destination string, contType string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	var stat os.FileInfo
	stat, err = os.Stat(source)
	if err != nil {
		return err
	}
	size := stat.Size()

	err = bucket.PutReader(destination, reader, size, contType, s3.PublicRead, s3.Options{})
	return err
}

func (u *Uploader) ListDirectory(directory string) ([]string, error) {
	var names []string
	prefix := directory
	delim := ""
	marker := ""
	max := 0
	resp, err := bucket.List(prefix, delim, marker, max)
	if err == nil {
		entries := resp.Contents
		for _, entry := range entries {
			name := filepath.Base(entry.Key)
			names = append(names, name)
		}
	}
	return names, err
}

// CreateDirectory does nothing since a directory does not need to be created on S3
// Directories are virtual, and defined by the path of the object
func (u *Uploader) CreateDirectory(path string) error {
	return nil
}

func Initialize(accessKey string, secretKey string, bucketName string) {
	bucket = retrieveBucket(accessKey, secretKey, bucketName)
}

func retrieveBucket(accessKey, secretKey, bucketName string) *s3.Bucket {
	client := s3.S3{
		Auth:   aws.Auth{AccessKey: accessKey, SecretKey: secretKey},
		Region: aws.USEast,
		AttemptStrategy: aws.AttemptStrategy{
			Min:   2,
			Total: 5 * time.Second,
			Delay: 200 * time.Millisecond,
		},
		ConnectTimeout: 2 * time.Second,
		ReadTimeout:    8 * time.Second,
		RequestTimeout: 10 * time.Second,
	}
	return client.Bucket(bucketName)
}
