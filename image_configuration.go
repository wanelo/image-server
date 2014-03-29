package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ImageConfiguration struct {
	id        string
	width     int
	height    int
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

func (ic *ImageConfiguration) ImageName() string {
	if ic.width == 0 && ic.height == 0 {
		return "full_size." + ic.format
	} else if ic.height == 0 {
		return fmt.Sprintf("w%d.%s", ic.width, ic.format)
	} else {
		return fmt.Sprintf("%dx%d.%s", ic.width, ic.height, ic.format)
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
	return ic.LocalDestinationDirectory() + "/" + ic.ImageName()
}

func (ic *ImageConfiguration) MantaOriginalImagePath() string {
	return serverConfiguration.MantaBasePath + "/" + ic.ImageDirectory() + "/original"
}

func (ic *ImageConfiguration) MantaResizedImagePath() string {
	return serverConfiguration.MantaBasePath + "/" + ic.ImageDirectory() + "/" + ic.ImageName()
}

func newImageConfiguration(r *http.Request) *ImageConfiguration {
	params := mux.Vars(r)
	qs := r.URL.Query()
	width, _ := strconv.Atoi(params["width"])
	height, _ := strconv.Atoi(params["height"])

	return &ImageConfiguration{
		model:     params["model"],
		imageType: params["imageType"],
		id:        params["id"],
		format:    params["format"],
		source:    qs.Get("source"),
		quality:   serverConfiguration.DefaultCompression,
		width:     width,
		height:    height,
	}
}
