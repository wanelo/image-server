package core

import (
	"fmt"
	"mime"
)

// ImageConfiguration struct
// Properties used to generate new image
type ImageConfiguration struct {
	// ServerConfiguration *ServerConfiguration
	ID        string
	Width     int
	Height    int
	Filename  string
	Format    string
	Source    string
	Quality   uint
	Namespace string
}

// ToContentType returns the content type based on the image format
func (ic *ImageConfiguration) ToContentType() string {
	ext := fmt.Sprintf(".%s", ic.Format)
	return mime.TypeByExtension(ext)
}
