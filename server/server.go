package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/fetcher"
	"github.com/wanelo/image-server/info"
	"github.com/wanelo/image-server/uploader"
)

func InitializeRouter(sc *core.ServerConfiguration, port string) {
	log.Println("starting server on http://0.0.0.0:" + port)

	r := render.New(render.Options{
		IndentJSON: true,
	})

	router := mux.NewRouter()
	router.HandleFunc("/{namespace:[a-z0-9]+}", func(wr http.ResponseWriter, req *http.Request) {
		NewImageHandler(wr, req, sc, r)
	}).Methods("POST").Name("newImage")

	// n := negroni.New()
	n := negroni.Classic()
	n.UseHandler(router)

	n.Run(":" + port)
}

func NewImageHandler(w http.ResponseWriter, req *http.Request, sc *core.ServerConfiguration, r *render.Render) {
	qs := req.URL.Query()
	vars := mux.Vars(req)
	errorStr := ""

	source := qs.Get("source")
	namespace := vars["namespace"]

	log.Printf("Processing request for: %s", source)

	f := fetcher.NewSourceFetcher(sc.Adapters.Paths, sc.Adapters.Fetcher)
	fc := f.Channels
	err, imageDetails := f.Fetch(source, namespace)
	var json map[string]string

	if err != nil {
		errorStr = fmt.Sprintf("%s", err)
		// r.JSON(w, http.StatusOK, json)
		json = map[string]string{
			"error": errorStr,
		}
		r.JSON(w, http.StatusOK, json)
		return
	}

	hash := imageDetails.Hash

	// go func() {
	select {
	case localOriginalPath := <-fc.DownloadComplete:
		uploader := &uploader.Uploader{sc.RemoteBasePath}
		err := uploader.CreateDirectory(sc.Adapters.Paths.RemoteImageDirectory(namespace, hash))
		if err != nil {
			log.Printf("Manta::sentToManta unable to create directory %s", sc.RemoteBasePath)
			return
		}

		destination := sc.Adapters.Paths.RemoteOriginalPath(namespace, hash)

		go sc.Adapters.Logger.OriginalDownloaded(localOriginalPath, destination)
		go func() {
			localInfoPath := sc.Adapters.Paths.LocalInfoPath(namespace, hash)
			remoteInfoPath := sc.Adapters.Paths.RemoteInfoPath(namespace, hash)

			err := info.SaveImageDetail(imageDetails, localInfoPath)
			if err != nil {
				log.Println(err)
			}

			// upload info
			err = uploader.Upload(localInfoPath, remoteInfoPath)
			if err != nil {
				log.Println(err)
			}
		}()

		// upload original image
		err = uploader.Upload(localOriginalPath, destination)
		if err != nil {
			log.Println(err)
		}

	case <-fc.DownloadFailed:
		go sc.Adapters.Logger.OriginalDownloadFailed(source)
	case <-fc.SkippedDownload:
		go sc.Adapters.Logger.OriginalDownloadSkipped(source)
	}
	// }()

	json = map[string]string{
		"error":  errorStr,
		"hash":   hash,
		"height": fmt.Sprintf("%v", imageDetails.Height),
		"width":  fmt.Sprintf("%v", imageDetails.Width),
	}

	r.JSON(w, http.StatusOK, json)
}
