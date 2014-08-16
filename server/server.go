package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"code.google.com/p/go-uuid/uuid"
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

	router.HandleFunc("/{namespace}/batch", func(wr http.ResponseWriter, req *http.Request) {
		BatchHandler(wr, req, sc)
	}).Methods("POST").Name("batch")

	// n := negroni.New()
	n := negroni.Classic()
	n.UseHandler(router)

	n.Run(":" + port)
}

func BatchHandler(w http.ResponseWriter, req *http.Request, sc *core.ServerConfiguration) {
	name := uuid.NewRandom()
	// func Post(url string, bodyType string, body io.Reader) (resp *Response, err error)
	// uploader := uploader.DefaultUploader(sc.RemoteBasePath)

	body := &bytes.Buffer{}
	_, err := body.ReadFrom(req.Body)
	defer req.Body.Close()

	if err == nil {
		os.MkdirAll("tmp", 0700)
		localPath := fmt.Sprintf("tmp/%s", name)
		log.Println("About to create", localPath)
		out, err := os.Create(localPath)
		if err != nil {
			panic(err)
		}
		defer out.Close()

		_, err = io.Copy(out, body)
		if err != nil {
			fmt.Fprintln(w, err)
		}

		// http.Post(fmt.Sprintf("http://localhost:3333/%s", name), "text/plain", body)

		fmt.Println(name)
	} else {
		log.Println(err)
	}

	// log.Println(fmt.Sprintf("http://localhost:3333/%s", name))

	// reader := bufio.NewReader(req.Body)
	//
	// var err error
	// eof := false
	// for !eof {
	// 	var line string
	// 	line, err = reader.ReadString('\n')
	// 	if err == io.EOF {
	// 		err = nil
	// 		eof = true
	// 	} else if err != nil {
	// 		return
	// 	}
	//
	// 	fmt.Println(line)
	// }

}

func NewImageHandler(w http.ResponseWriter, req *http.Request, sc *core.ServerConfiguration) {
	r := render.New(render.Options{
		IndentJSON: true,
	})

	qs := req.URL.Query()
	vars := mux.Vars(req)

	source := qs.Get("source")
	namespace := vars["namespace"]

	log.Printf("Processing request for: %s", source)

	f := fetcher.NewSourceFetcher(sc.Adapters.Paths)

	imageDetails, downloaded, err := f.Fetch(source, namespace)
	var json map[string]string

	if err != nil {
		errorHandlerJSON(err, w, r, http.StatusNotFound)
		return
	}

	hash := imageDetails.Hash

	if downloaded {
		localOriginalPath := f.Paths.LocalOriginalPath(namespace, hash)
		uploader := uploader.DefaultUploader(sc.RemoteBasePath)
		err := uploader.CreateDirectory(sc.Adapters.Paths.RemoteImageDirectory(namespace, hash))
		if err != nil {
			log.Printf("Manta::sentToManta unable to create directory %s", sc.RemoteBasePath)
			errorHandlerJSON(err, w, r, http.StatusInternalServerError)
			return
		}

		destination := sc.Adapters.Paths.RemoteOriginalPath(namespace, hash)

		go sc.Adapters.Logger.OriginalDownloaded(localOriginalPath, destination)

		localInfoPath := sc.Adapters.Paths.LocalInfoPath(namespace, hash)
		remoteInfoPath := sc.Adapters.Paths.RemoteInfoPath(namespace, hash)

		err = info.SaveImageDetail(imageDetails, localInfoPath)
		if err != nil {
			log.Println(err)
		}

		// upload info
		err = uploader.Upload(localInfoPath, remoteInfoPath)
		if err != nil {
			log.Println(err)
		}

		// upload original image
		err = uploader.Upload(localOriginalPath, destination)
		if err != nil {
			log.Println(err)
			errorHandlerJSON(err, w, r, http.StatusInternalServerError)
			return
		}
	}

	json = map[string]string{
		"hash":   hash,
		"height": fmt.Sprintf("%v", imageDetails.Height),
		"width":  fmt.Sprintf("%v", imageDetails.Width),
	}

	if err != nil {
		json["error"] = fmt.Sprintf("%s", err)
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
		uploader := uploader.DefaultUploader(sc.RemoteBasePath)
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

func errorHandlerJSON(err error, w http.ResponseWriter, r *render.Render, status int) {
	json := map[string]string{
		"error": fmt.Sprintf("%s", err),
	}
	r.JSON(w, status, json)
}

func errorHandler(err error, w http.ResponseWriter, r *http.Request, status int, sc *core.ServerConfiguration, ic *core.ImageConfiguration) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 image not available. ", err)
	}
}
