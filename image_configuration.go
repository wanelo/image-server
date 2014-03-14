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
	quality   int
	model     string
	imageType string
}

func (ic *ImageConfiguration) RemoteImageUrl() string {
	if ic.source != "" {
		return ic.source
	} else {
		return "http://cdn-s3-2.wanelo.com/" + ic.model + "/" + ic.imageType + "/" + ic.id + "/original.jpg"
	}
}

func (ic *ImageConfiguration) OriginalImagePath() string {
	return "public/" + ic.model + "/" + ic.imageType + "/" + ic.id + "/original"
}

func (ic *ImageConfiguration) ResizedImagePath() string {
	if ic.width == 0 && ic.height == 0 {
		return "public/" + ic.model + "/" + ic.imageType + "/" + ic.id + "/full_size." + ic.format
	} else if ic.height == 0 {
		return fmt.Sprintf("public/%s/%s/%s/w%d.%s", ic.model, ic.imageType, ic.id, ic.width, ic.format)
	} else {
		return fmt.Sprintf("public/%s/%s/%s/%dx%d.%s", ic.model, ic.imageType, ic.id, ic.width, ic.height, ic.format)
	}
}

func (ic *ImageConfiguration) MagickInfo() *magick.Info {
	info := magick.NewInfo()
	info.SetQuality(75)
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

	return ic
}
