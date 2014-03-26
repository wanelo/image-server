package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rainycape/magick"
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
		return serverConfiguration.SourceDomain + "/" + ic.model + "/" + ic.imageType + "/" + ic.id + "/original.jpg"
	}
}

func (ic *ImageConfiguration) OriginalImagePath() string {
	return "public/" + ic.model + "/" + ic.imageType + "/" + ic.id + "/original"
}

func (ic *ImageConfiguration) DestinationDirectory() string {
	return fmt.Sprintf("public/%s/%s/%s", ic.model, ic.imageType, ic.id)
}

func (ic *ImageConfiguration) ResizedImagePath() string {
	dir := ic.DestinationDirectory()

	if ic.width == 0 && ic.height == 0 {
		return dir + "/full_size." + ic.format
	} else if ic.height == 0 {
		return fmt.Sprintf("%s/w%d.%s", dir, ic.width, ic.format)
	} else {
		return fmt.Sprintf("%s/%dx%d.%s", dir, ic.width, ic.height, ic.format)
	}
}

func (ic *ImageConfiguration) MagickInfo() *magick.Info {
	info := magick.NewInfo()
	info.SetQuality(ic.quality)
	info.SetFormat(ic.format)
	return info
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
		quality:   75,
		width:     width,
		height:    height,
	}
}
