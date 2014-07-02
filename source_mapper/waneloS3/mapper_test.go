package mapper

import (
	"testing"

	"github.com/wanelo/image-server/core"
)

func TestRemoteImageURL(t *testing.T) {
	sc := &core.ServerConfiguration{
		SourceDomain:  "http://example.com",
		MaximumWidth:  1000,
		LocalBasePath: "public",
	}

	ic := &core.ImageConfiguration{
		Namespace:           "p",
		ID:                  "00ofrA",
	}

	mapper := SourceMapper{sc}

	expected := "http://example.com/product/image/12077300/original.jpg"
	remoteURL := mapper.RemoteImageURL(ic)

	if expected != remoteURL {
		t.Errorf("expected\n%v got\n%v", expected, remoteURL)
	}
}
