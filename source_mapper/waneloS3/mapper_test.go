package mapper

import (
	"testing"

	"github.com/wanelo/image-server/core"
)

func TestRemoteImageURL(t *testing.T) {
	mappings := make(map[string]string)
	mappings["p"] = "test/image"

	sc := &core.ServerConfiguration{
		SourceDomain:  "http://example.com",
		MaximumWidth:  1000,
		LocalBasePath: "public",
	}

	mc := &core.MapperConfiguration{mappings}

	ic := &core.ImageConfiguration{
		ServerConfiguration: sc,
		Namespace:           "p",
		ID:                  "00ofrA",
	}

	mapper := SourceMapper{mc}

	expected := "http://example.com/test/image/12077300/original.jpg"
	remoteURL := mapper.RemoteImageURL(ic)

	if expected != remoteURL {
		t.Errorf("expected\n%v got\n%v", expected, remoteURL)
	}
}
