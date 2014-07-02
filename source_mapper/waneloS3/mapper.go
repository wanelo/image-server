package mapper

import (
	"fmt"

	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/encoders/base62"
)

type SourceMapper struct{
	ServerConfiguration *core.ServerConfiguration
}

// RemoteImageURL returns a URL string for original image
func (m *SourceMapper) RemoteImageURL(ic *core.ImageConfiguration) string {
	if ic.Source != "" {
		return ic.Source
	}

	url := m.ServerConfiguration.SourceDomain + "/" + m.imageDirectory(ic) + "/original.jpg"
	return url
}

func (m *SourceMapper) imageDirectory(ic *core.ImageConfiguration) string {
	id := base62.Decode(ic.ID)
	return fmt.Sprintf("%s/%d", namespaceMapping(ic.Namespace), id)
}

func namespaceMapping(namespace string) string {
	switch {
	case namespace == "p":
		return "product/image"
	}
	return namespace
}
