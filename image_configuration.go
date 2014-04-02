package main

import "fmt"

type ImageConfiguration struct {
	id        string
	width     int
	height    int
	filename  string
	format    string
	source    string
	quality   uint
	model     string
	imageType string
}

func (ic *ImageConfiguration) RemoteImageUrl() string {
	if ic.source != "" {
		return ic.source
	} else {
		return serverConfiguration.SourceDomain + "/" + ic.ImageDirectory() + "/original.jpg"
	}
}

func (ic *ImageConfiguration) ImageDirectory() string {
	return fmt.Sprintf("%s/%s/%s", ic.model, ic.imageType, ic.id)
}

func (ic *ImageConfiguration) LocalDestinationDirectory() string {
	return "public/" + ic.ImageDirectory()
}

func (ic *ImageConfiguration) LocalOriginalImagePath() string {
	return ic.LocalDestinationDirectory() + "/original"
}

func (ic *ImageConfiguration) LocalResizedImagePath() string {
	return ic.LocalDestinationDirectory() + "/" + ic.filename
}

func (ic *ImageConfiguration) MantaOriginalImagePath() string {
	return serverConfiguration.MantaBasePath + "/" + ic.ImageDirectory() + "/original"
}

func (ic *ImageConfiguration) MantaResizedImagePath() string {
	return serverConfiguration.MantaBasePath + "/" + ic.ImageDirectory() + "/" + ic.filename
}
