package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rainycape/magick"
	"net/http"
	"strconv"
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

func buildImageConfiguration(r *http.Request) *ImageConfiguration {
	ic := new(ImageConfiguration)
	params := mux.Vars(r)
	qs := r.URL.Query()

	ic.model = params["model"]
	ic.imageType = params["imageType"]
	ic.id = params["id"]
	ic.width, _ = strconv.Atoi(params["width"])
	ic.height, _ = strconv.Atoi(params["height"])
	ic.format = params["format"]
	ic.source = qs.Get("source")
	ic.quality = 75

	return ic
}
