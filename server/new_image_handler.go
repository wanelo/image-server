package server

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/info"
	"github.com/wanelo/image-server/request"
	"github.com/wanelo/image-server/uploader"
)

// NewImageHandler handles posting new images
func NewImageHandler(w http.ResponseWriter, req *http.Request, sc *core.ServerConfiguration) {
	qs := req.URL.Query()
	vars := mux.Vars(req)
	source := qs.Get("source")
	namespace := vars["namespace"]

	request := &request.Request{
		ServerConfiguration: sc,
		Namespace:           namespace,
		Outputs:             strings.Split(qs.Get("outputs"), ","),
		Uploader:            uploader.DefaultUploader(sc),
		Paths:               sc.Adapters.Paths,
		SourceURL:           source,
	}

	imageDetails, err := request.Create()
	if err != nil {
		errorHandlerJSON(err, w, http.StatusNotFound)
		return
	}

	renderImageDetails(w, imageDetails)
}

func renderImageDetails(w http.ResponseWriter, imageDetails *info.ImageDetails) {
	r := render.New(render.Options{
		IndentJSON: true,
	})

	r.JSON(w, http.StatusOK, imageDetails)
}
