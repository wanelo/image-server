package s3_test

import (
	"log"
	"os"
	"testing"

	. "github.com/image-server/image-server/test"
	"github.com/image-server/image-server/uploader/s3"
)

func TestItemToHash(t *testing.T) {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		return
	}

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_KEY")
	bucketName := os.Getenv("AWS_BUCKET")

	s3.Initialize(accessKey, secretKey, bucketName)

	uploader := s3.Uploader{
		BaseDir: os.Getenv("IMG_REMOTE_BASE_PATH"),
	}

	existing, err := uploader.ListDirectory("p/543/47c/442/1c41f9467a3f5afed64943b")
	Ok(t, err)
	log.Println(existing)

	// Equals(t, "6ad5544baa6f5e852e1af26f8c2e45db", image.ToHash())
}
