package server

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/image-server/image-server/core"
	"github.com/image-server/image-server/request"
	"github.com/image-server/image-server/uploader"
	"github.com/unrolled/render"
)

type NewFileResponse struct {
}

// NewFileHandler handles posting new files
func NewFileHandler(w http.ResponseWriter, req *http.Request, sc *core.ServerConfiguration) {
	qs := req.URL.Query()
	vars := mux.Vars(req)
	sourceURL := qs.Get("source")
	namespace := vars["namespace"]
	filename := vars["filename"]
	hash := varsToHash(vars)

	request := &request.Request{
		ServerConfiguration: sc,
		Namespace:           namespace,
		Uploader:            uploader.DefaultUploader(sc),
		Paths:               sc.Adapters.Paths,
		Hash:                hash,
		SourceURL:           sourceURL,
		SourceData:          req.Body,
	}

	err := request.UploadFile(filename)
	if err != nil {
		glog.Error("Failed to upload file from ", sourceURL, " - ", err)
		errorHandlerJSON(err, w, http.StatusNotFound)
		return
	}

	response := &NewFileResponse{}

	renderUploadDetails(w, response)
}

func renderUploadDetails(w http.ResponseWriter, res *NewFileResponse) {
	r := render.New(render.Options{
		IndentJSON: true,
	})

	r.JSON(w, http.StatusOK, res)
}
