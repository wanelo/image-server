package s3

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

// Uploader for S3
type Uploader struct {
	AccessKey  string
	SecretKey  string
	Token      string
	BucketName string
	BaseDir    string
}

// Upload copies a file int a bucket in S3
func (u *Uploader) Upload(source string, destination string, contType string) error {
	bucket := u.retrieveBucket()
	reader, err := os.Open(source)

	if err != nil {
		return err
	}

	if contType == "" {
		data, err := ioutil.ReadFile(source)
		if err != nil {
			return err
		}
		contType = http.DetectContentType(data)
	}

	var stat os.FileInfo
	stat, err = os.Stat(source)
	if err != nil {
		return err
	}
	size := stat.Size()

	err = bucket.PutReader(destination, reader, size, contType, s3.PublicRead)
	return err
}

func (u *Uploader) ListDirectory(directory string) ([]string, error) {
	var names []string
	bucket := u.retrieveBucket()
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

// Initialize does nothing
func (u *Uploader) Initialize() error {
	return nil
}

func (u *Uploader) retrieveBucket() *s3.Bucket {
	auth := aws.Auth{AccessKey: u.AccessKey, SecretKey: u.SecretKey, Token: u.Token}
	client := s3.New(auth, aws.USEast)
	return client.Bucket(u.BucketName)
}
