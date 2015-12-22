package core_test

import (
	"testing"

	"github.com/image-server/image-server/core"
	. "github.com/image-server/image-server/test"
)

func TestToContentTypeForJPEG(t *testing.T) {
	ic := &core.ImageConfiguration{Format: "jpeg"}

	Equals(t, "image/jpeg", ic.ToContentType())
}

func TestToContentTypeForJpg(t *testing.T) {
	ic := &core.ImageConfiguration{Format: "jpg"}

	Equals(t, "image/jpeg", ic.ToContentType())
}

func TestToContentTypeForWebp(t *testing.T) {
	ic := &core.ImageConfiguration{Format: "webp"}

	Equals(t, "image/webp", ic.ToContentType())
}

func TestToContentTypeForGIF(t *testing.T) {
	ic := &core.ImageConfiguration{Format: "gif"}

	Equals(t, "image/gif", ic.ToContentType())
}
