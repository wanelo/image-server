package info_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/wanelo/image-server/info"
	. "github.com/wanelo/image-server/test"
)

func TestImageHash(t *testing.T) {
	i := info.Info{"../test/images/a.jpg"}
	expectedHash := "31e8b3187a9f63f26d58c88bf09a7bbd"
	hash, err := i.FileHash()
	Ok(t, err)
	Equals(t, expectedHash, hash)
}

func TestImageDetailsOnJPEG(t *testing.T) {
	i := info.Info{"../test/images/a.jpg"}
	imageDetails, err := i.ImageDetails()
	expectedHash := "31e8b3187a9f63f26d58c88bf09a7bbd"
	Ok(t, err)
	Equals(t, expectedHash, imageDetails.Hash)
	Equals(t, 496, imageDetails.Height)
	Equals(t, 574, imageDetails.Width)
	Equals(t, "image/jpeg", imageDetails.ContentType)
}

func TestImageDetailsOnPNG(t *testing.T) {
	i := info.Info{"../test/images/a.png"}
	imageDetails, err := i.ImageDetails()
	expectedHash := "117813b6a51e74c77d0fc7d5de510f42"
	Ok(t, err)
	Equals(t, expectedHash, imageDetails.Hash)
	Equals(t, imageDetails.Height, 496)
	Equals(t, imageDetails.Width, 574)
	Equals(t, "image/png", imageDetails.ContentType)
}

func TestImageDetailsOnWEBP(t *testing.T) {
	i := info.Info{"../test/images/a.webp"}
	imageDetails, err := i.ImageDetails()
	expectedHash := "2a9d1753531a2c060c002a97b983854c"
	Ok(t, err)
	Equals(t, expectedHash, imageDetails.Hash)
	Equals(t, imageDetails.Height, 496)
	Equals(t, imageDetails.Width, 574)
	Equals(t, "image/webp", imageDetails.ContentType)
}

func TestImageDetailsOnWEBPWihtoutExtension(t *testing.T) {
	i := info.Info{"../test/images/webp_without_ext"}
	imageDetails, err := i.ImageDetails()
	expectedHash := "2a9d1753531a2c060c002a97b983854c"
	Ok(t, err)
	Equals(t, expectedHash, imageDetails.Hash)
	Equals(t, imageDetails.Height, 496)
	Equals(t, imageDetails.Width, 574)
	Equals(t, "image/webp", imageDetails.ContentType)
}

func TestImageDetailsToJSON(t *testing.T) {
	d := &info.ImageDetails{"THISISAHASH", 10, 20, "image/jpeg"}
	json, err := info.ImageDetailsToJSON(d)
	expected := "{\"hash\":\"THISISAHASH\",\"height\":10,\"width\":20,\"content_type\":\"image/jpeg\"}"
	Ok(t, err)
	Equals(t, expected, json)
}

func TestSaveImageDetail(t *testing.T) {
	path := "../test/test-image-detail.json"
	d := &info.ImageDetails{"THISISAHASH", 10, 20, "image/jpeg"}
	info.SaveImageDetail(d, path)

	fileBuffer, err := ioutil.ReadFile(path)
	Ok(t, err)

	expected := "{\"hash\":\"THISISAHASH\",\"height\":10,\"width\":20,\"content_type\":\"image/jpeg\"}"
	fileContents := string(fileBuffer)
	Equals(t, expected, fileContents)

	os.Remove(path)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}
