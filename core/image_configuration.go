package core

import "fmt"

// ImageConfiguration struct
// Properties used to generate new image
type ImageConfiguration struct {
	ServerConfiguration *ServerConfiguration
	ID                  string
	Width               int
	Height              int
	Filename            string
	Format              string
	Source              string
	Quality             uint
	Model               string
	ImageType           string
}

// RemoteImageURL returns a URL string for original image
func (ic *ImageConfiguration) RemoteImageURL() string {
	if ic.Source != "" {
		return ic.Source
	}
	return ic.ServerConfiguration.SourceDomain + "/" + ic.ImageDirectory() + "/original.jpg"
}

func (ic *ImageConfiguration) ImageDirectory() string {
	return fmt.Sprintf("%s/%s/%s", ic.Model, ic.ImageType, ic.ID)
}

func (ic *ImageConfiguration) LocalDestinationDirectory() string {
	return ic.ServerConfiguration.LocalBasePath + "/" + ic.ImageDirectory()
}

func (ic *ImageConfiguration) LocalOriginalImagePath() string {
	return ic.LocalDestinationDirectory() + "/original"
}

func (ic *ImageConfiguration) LocalResizedImagePath() string {
	return ic.LocalDestinationDirectory() + "/" + ic.Filename
}

func (sc *ServerConfiguration) MantaResizedImagePath(ic *ImageConfiguration) string {
	return sc.MantaBasePath + "/" + ic.ImageDirectory() + "/" + ic.Filename
}
