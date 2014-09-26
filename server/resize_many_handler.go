package server

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/uploader"

	"github.com/wanelo/image-server/request"
)

// ResizeManyHandler asumes the original image is either stores locally or on the remote server
// It returns status code 200 and no content
// A listing will be requested to the uploader to determine what images are missing, and only
// Images not already processed will be generated and uploaded
func ResizeManyHandler(w http.ResponseWriter, req *http.Request, sc *core.ServerConfiguration) {
	qs := req.URL.Query()
	vars := mux.Vars(req)

	ir := request.Request{
		ServerConfiguration: sc,
		Namespace:           vars["namespace"],
		Outputs:             strings.Split(qs.Get("outputs"), ","),
		Uploader:            uploader.DefaultUploader(sc),
		Paths:               sc.Adapters.Paths,
		Hash:                varsToHash(vars),
	}

	err := ir.ProcessMultiple()
	if err != nil {
		errorHandlerJSON(err, w, http.StatusInternalServerError)
		return
	}
}
