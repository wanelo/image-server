package info

import "testing"

func assert(t *testing.T, value int, expected int) {
  if value != expected {
    t.Errorf("expected %v to be %v", value, expected)
  }
}

func TestImageHash(t *testing.T) {
	i := Info{"../test/images/a.jpg"}
  expectedHash := "31e8b3187a9f63f26d58c88bf09a7bbd"
	if i.FileHash() != expectedHash {
		t.Errorf("expected %v to be %v", i.FileHash(), expectedHash)
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
