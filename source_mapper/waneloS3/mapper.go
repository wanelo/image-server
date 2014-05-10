package mapper

import (
	"fmt"

	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/encoders/base62"
)

type SourceMapper struct {
	ServerConfiguration *core.ServerConfiguration
}

// RemoteImageURL returns a URL string for original image
func (m *SourceMapper) RemoteImageURL(ic *core.ImageConfiguration) string {
	if ic.Source != "" {
		return ic.Source
	}
	url := ic.ServerConfiguration.SourceDomain + "/" + m.imageDirectory(ic) + "/original.jpg"
	return url
}

func (m *SourceMapper) imageDirectory(ic *core.ImageConfiguration) string {
	id := base62.Decode("ofrA")
	// fmt.Printf("Decoded %s to %d", ic.ID, id)
	return fmt.Sprintf("%s/%d", m.ServerConfiguration.NamespaceMappings[ic.Namespace], id)
}
