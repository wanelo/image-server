package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/fetcher"
	"github.com/wanelo/image-server/info"
	"github.com/wanelo/image-server/uploader"
)

// NewImageHandler handles posting new images
func NewImageHandler(w http.ResponseWriter, req *http.Request, sc *core.ServerConfiguration) {
	qs := req.URL.Query()
	vars := mux.Vars(req)
	source := qs.Get("source")
	namespace := vars["namespace"]
	f := fetcher.NewSourceFetcher(sc.Adapters.Paths)
	imageDetails, downloaded, err := f.Fetch(source, namespace)

	if err != nil {
		errorHandlerJSON(err, w, http.StatusNotFound)
		return
	}

	localOriginalPath := f.Paths.LocalOriginalPath(namespace, imageDetails.Hash)
	if downloaded {
		err := uploadOriginal(sc, namespace, imageDetails, localOriginalPath)
		if err != nil {
			errorHandlerJSON(err, w, http.StatusInternalServerError)
			return
		}
	}

	outputs := strings.Split(qs.Get("outputs"), ",")
	err = processAndUploadFromOutputs(sc, localOriginalPath, namespace, imageDetails.Hash, outputs)
	if err != nil {
		errorHandlerJSON(err, w, http.StatusNotFound)
		return
	}

	renderImageDetails(w, imageDetails)
}

func uploadOriginal(sc *core.ServerConfiguration, namespace string, imageDetails *info.ImageDetails, localOriginalPath string) error {
	uploader := uploader.DefaultUploader(sc)
	err := uploader.CreateDirectory(sc.Adapters.Paths.RemoteImageDirectory(namespace, imageDetails.Hash))
	if err != nil {
		return err
	}

	destination := sc.Adapters.Paths.RemoteOriginalPath(namespace, imageDetails.Hash)

	go sc.Adapters.Logger.OriginalDownloaded(localOriginalPath, destination)

	uploadImageDetails(sc, namespace, imageDetails, uploader)

	// upload original image
	err = uploader.Upload(localOriginalPath, destination, "")
	if err != nil {
		return err
	}
	return nil
}

// uploadImageDetails uploads info.json
func uploadImageDetails(sc *core.ServerConfiguration, namespace string, imageDetails *info.ImageDetails, uploader *uploader.Uploader) {
	localInfoPath := sc.Adapters.Paths.LocalInfoPath(namespace, imageDetails.Hash)
	remoteInfoPath := sc.Adapters.Paths.RemoteInfoPath(namespace, imageDetails.Hash)

	err := info.SaveImageDetail(imageDetails, localInfoPath)
	if err != nil {
		log.Println(err)
		return
	}

	// upload info
	err = uploader.Upload(localInfoPath, remoteInfoPath, "application/json")
	if err != nil {
		log.Println(err)
	}
}

func renderImageDetails(w http.ResponseWriter, imageDetails *info.ImageDetails) {
	r := render.New(render.Options{
		IndentJSON: true,
	})

	r.JSON(w, http.StatusOK, imageDetails)
}
