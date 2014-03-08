package main

import (
	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
	"image/jpeg"
	"io"
	/*	"io/ioutil"*/
	"log"
	"net/http"
	"os"
	/*	"strconv"*/)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/product/image/{id:[0-9]+}/{width:[0-9]+}x{height:[0-9]+}.jpg", imageHandler).Methods("GET")
	http.Handle("/", r)
	log.Println("Listening on port 3000...")
	http.ListenAndServe(":3000", nil)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	width := params["width"]
	height := params["height"]

	resp, err := http.Get("http://cdn-s3-2.wanelo.com/product/image/" + id + "/original.jpg")
	if err != nil {
		// panic?
	}
	defer resp.Body.Close()
	out, err := os.Create("public/" + id)
	defer out.Close()

	/*	imgBody := resp.Body*/
	io.Copy(out, resp.Body)

	file, err := os.Open("public/" + id)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	outResized, err := os.Create("public/" + id + "_" + width + "x" + height + ".jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outResized.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	/*	widthInt, _ := strconv.ParseInt(width, 32, 2)*/
	m := resize.Resize(100, 100, img, resize.Lanczos3)

	// write new image to file
	jpeg.Encode(outResized, m, nil)

	/*	w.ServeContent(w, req, name, modtime, sizeFunc, content)*/

	w.Write([]byte("Hello " + id + " " + width + " " + height))
}
