package info

import (
	"os"
	"testing"
)

func assert(t *testing.T, value int, expected int) {
	if value != expected {
		t.Errorf("expected %v to be %v", value, expected)
	}
}

func TestImageHash(t *testing.T) {
	i := Info{"../test/images/a.jpg"}
	expectedHash := "31e8b3187a9f63f26d58c88bf09a7bbd"
	hash, _ := i.FileHash()
	if hash != expectedHash {
		t.Errorf("expected %v to be %v", hash, expectedHash)
	}
}

func TestImageDetailsOnJPEG(t *testing.T) {
	i := Info{"../test/images/a.jpg"}
	imageDetails, _ := i.ImageDetails()
	expectedHash := "31e8b3187a9f63f26d58c88bf09a7bbd"
	if imageDetails.Hash != expectedHash {
		t.Errorf("expected %v to be %v", imageDetails.Hash, expectedHash)
	}

	assert(t, imageDetails.Height, 496)
	assert(t, imageDetails.Width, 574)
}

func TestImageDetailsOnPNG(t *testing.T) {
	i := Info{"../test/images/a.png"}
	imageDetails, _ := i.ImageDetails()
	expectedHash := "117813b6a51e74c77d0fc7d5de510f42"
	if imageDetails.Hash != expectedHash {
		t.Errorf("expected %v to be %v", imageDetails.Hash, expectedHash)
	}

	assert(t, imageDetails.Height, 496)
	assert(t, imageDetails.Width, 574)
}

// func TestImageDetailsOnWEBP(t *testing.T) {
//   i := Info{"../test/images/a.webp"}
//   imageDetails, _ := i.ImageDetails()
//   expectedHash := "117813b6a51e74c77d0fc7d5de510f42"
//   if imageDetails.Hash != expectedHash {
//     t.Errorf("expected %v to be %v", imageDetails.Hash, expectedHash)
//   }
//
//   assert(t, imageDetails.Height, 496)
//   assert(t, imageDetails.Width, 574)
// }

func TestImageDetailsToJSON(t *testing.T) {
	d := &ImageDetails{"THISISAHASH", 10, 20}
	json, _ := ImageDetailsToJSON(d)
	expected := "{\"hash\":\"THISISAHASH\",\"height\":10,\"width\":20}"

	if expected != json {
		t.Errorf("expected %v to be %v", json, expected)
	}
}

func TestSaveImageDetail(t *testing.T) {
	path := "../test/test-image-detail.json"
	d := &ImageDetails{"THISISAHASH", 10, 20}
	SaveImageDetail(d, path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file to exist on %v", path)
		return
	}

	os.Remove(path)
}
