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
	"github.com/wanelo/image-server/parser"
	"github.com/wanelo/image-server/processor"
	"github.com/wanelo/image-server/uploader"
)

func InitializeRouter(sc *core.ServerConfiguration, port string) {
	log.Println("starting server on http://0.0.0.0:" + port)

	router := mux.NewRouter()
	router.HandleFunc("/{namespace:[a-z0-9]+}", func(wr http.ResponseWriter, req *http.Request) {
		NewImageHandler(wr, req, sc)
	}).Methods("POST").Name("newImage")

	router.HandleFunc("/{namespace:[a-z0-9]+}/{id1:[a-f0-9]{3}}/{id2:[a-f0-9]{3}}/{id3:[a-f0-9]{3}}/{id4:[a-f0-9]{23}}/{filename}", func(wr http.ResponseWriter, req *http.Request) {
		ResizeHandler(wr, req, sc)
	}).Methods("GET").Name("resizeImage")

	// n := negroni.New()
	n := negroni.Classic()
	n.UseHandler(router)

	n.Run(":" + port)
}

func NewImageHandler(w http.ResponseWriter, req *http.Request, sc *core.ServerConfiguration) {
	r := render.New(render.Options{
		IndentJSON: true,
	})

	qs := req.URL.Query()
	vars := mux.Vars(req)
	errorStr := ""

	source := qs.Get("source")
	namespace := vars["namespace"]

	log.Printf("Processing request for: %s", source)

	f := fetcher.NewSourceFetcher(sc.Adapters.Paths)

	imageDetails, downloaded, err := f.Fetch(source, namespace)
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

	if downloaded {
		localOriginalPath := f.Paths.LocalOriginalPath(namespace, hash)
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
	}

	json = map[string]string{
		"error":  errorStr,
		"hash":   hash,
		"height": fmt.Sprintf("%v", imageDetails.Height),
		"width":  fmt.Sprintf("%v", imageDetails.Width),
	}

	r.JSON(w, http.StatusOK, json)
}

func ResizeHandler(w http.ResponseWriter, req *http.Request, sc *core.ServerConfiguration) {
	vars := mux.Vars(req)
	filename := vars["filename"]

	ic, err := parser.NameToConfiguration(sc, filename)
	if err != nil {
		errorHandler(err, w, req, http.StatusNotFound, sc, ic)
		return
	}

	namespace := vars["namespace"]
	id1 := vars["id1"]
	id2 := vars["id2"]
	id3 := vars["id3"]
	id4 := vars["id4"]
	hash := fmt.Sprintf("%s%s%s%s", id1, id2, id3, id4)

	ic.ID = hash
	ic.Namespace = namespace

	localPath := sc.Adapters.Paths.LocalImagePath(namespace, hash, filename)
	localOriginalPath := sc.Adapters.Paths.LocalOriginalPath(namespace, hash)

	// download original image
	remoteOriginalPath := sc.Adapters.Paths.RemoteOriginalURL(namespace, hash)
	log.Println(remoteOriginalPath)
	f := fetcher.NewUniqueFetcher(remoteOriginalPath, localOriginalPath)
	_, err = f.Fetch()
	if err != nil {
		errorHandler(err, w, req, http.StatusNotFound, sc, ic)
		return
	}

	// process image
	pchan := &processor.ProcessorChannels{
		ImageProcessed: make(chan *core.ImageConfiguration),
		Skipped:        make(chan string),
	}
	defer close(pchan.ImageProcessed)
	defer close(pchan.Skipped)

	p := processor.Processor{
		Source:             localOriginalPath,
		Destination:        localPath,
		ImageConfiguration: ic,
		Channels:           pchan,
	}

	resizedPath, err := p.CreateImage()

	if err != nil {
		errorHandler(err, w, req, http.StatusNotFound, sc, ic)
		return
	}

	select {
	case <-pchan.ImageProcessed:
		log.Println("about to upload to manta")
		uploader := &uploader.Uploader{sc.RemoteBasePath}
		remoteResizedPath := sc.Adapters.Paths.RemoteImagePath(namespace, hash, filename)
		err = uploader.Upload(localPath, remoteResizedPath)
		if err != nil {
			log.Println(err)
		}
	case path := <-pchan.Skipped:
		log.Printf("Skipped processing %s", path)
	}

	http.ServeFile(w, req, resizedPath)
}

func errorHandler(err error, w http.ResponseWriter, r *http.Request, status int, sc *core.ServerConfiguration, ic *core.ImageConfiguration) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 image not available. ", err)
	}
}
