package s3_test

import (
	"log"
	"os"
	"testing"

	. "github.com/image-server/image-server/test"
	"github.com/image-server/image-server/uploader/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
)

func TestItemToHash(t *testing.T) {
	bucketName := os.Getenv("AWS_BUCKET")
	regionName := os.Getenv("AWS_REGION")

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(regionName),
	}))
	Assert(t, sess != nil, "We need AWS access for integration tests")

	s3.Initialize(bucketName, regionName)

	uploader := s3.Uploader{}

	existing, err := uploader.ListDirectory("p/543/47c/442/1c41f9467a3f5afed64943b")
	Ok(t, err)
	log.Println(existing)

	// Equals(t, "6ad5544baa6f5e852e1af26f8c2e45db", image.ToHash())
}
