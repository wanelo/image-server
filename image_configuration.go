package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ImageConfiguration struct {
	id        string
	width     int
	height    int
	widthStr  string
	heightStr string
	format    string
	source    string
}

func (ic *ImageConfiguration) RemoteImageUrl() string {
	if ic.source != "" {
		return ic.source
	} else {
		return "http://cdn-s3-2.wanelo.com/product/image/" + ic.id + "/original.jpg"
	}
}

func (ic *ImageConfiguration) OriginalImagePath() string {
	return "public/" + ic.id
}

func (ic *ImageConfiguration) ResizedImagePath() string {
	if ic.height == 0 {
		return fmt.Sprintf("public/generated/%s_x%d.%s", ic.id, ic.width, ic.format)
	} else {
		return fmt.Sprintf("public/generated/%s_%dx%d.%s", ic.id, ic.width, ic.height, ic.format)
	}
}

func buildImageConfiguration(r *http.Request) *ImageConfiguration {
	ic := new(ImageConfiguration)
	params := mux.Vars(r)
	qs := r.URL.Query()

	ic.id = params["id"]
	ic.width, _ = strconv.Atoi(params["width"])
	ic.height, _ = strconv.Atoi(params["height"])
	ic.format = params["format"]
	ic.source = qs.Get("source")

	return ic
}
