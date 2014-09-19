package s3

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

type Uploader struct {
	BucketName string
	Bucket     s3.Bucket
}

// Upload copies a file int a bucket in S3
func (u *Uploader) Upload(source string, destination string) error {
	bucket, err := u.bucket()
	var data []byte

	// need real content type here
	contType := "image/jpeg"

	perm := s3.PublicRead
	data, err = pathToBytes(source)

	err = bucket.Put(destination, data, contType, perm)
	// log.Print(fmt.Sprintf("%T %+v", resp.Buckets[0], resp.Buckets[0]))
	return err
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

func (u *Uploader) bucket() (*s3.Bucket, error) {
	auth, err := aws.EnvAuth()

	if err != nil {
		return &s3.Bucket{}, err
	}

	client := s3.New(auth, aws.USEast)
	return client.Bucket(u.BucketName), nil
}

func pathToBytes(path string) ([]byte, error) {
	var data []byte
	return data, nil
}
