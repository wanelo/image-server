package main

import (
	"code.google.com/p/graphics-go/graphics"
	"github.com/gorilla/mux"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/product/image/{id:[0-9]+}/{width:[0-9]+}x{height:[0-9]+}.jpg", widthHeightHandler).Methods("GET")
	r.HandleFunc("/product/image/{id:[0-9]+}/{width:[0-9]+}.jpg", squareHandler).Methods("GET")
	http.Handle("/", r)
	log.Println("Listening on port 3000...")
	http.ListenAndServe(":3000", nil)
}

func downloadAndSaveOriginal(path string, productId string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		resp, err := http.Get("http://cdn-s3-2.wanelo.com/product/image/" + productId + "/original.jpg")
		if err != nil {
			// panic?
		}
		defer resp.Body.Close()
		out, err := os.Create(path)
		defer out.Close()

		/*	imgBody := resp.Body*/
		io.Copy(out, resp.Body)
	}
}

func createResizedImage(fullSizePath string, resizedPath string, width string, height string) {
	file, err := os.Open(fullSizePath)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(resizedPath); os.IsNotExist(err) {
		// Third way
		widthInt, _ := strconv.Atoi(width)
		heightInt, _ := strconv.Atoi(height)

		dst := image.NewRGBA(image.Rect(0, 0, widthInt, heightInt))
		graphics.Thumbnail(dst, img)

		toimg, err := os.Create(resizedPath)
		if err != nil {
			log.Fatal(err)
		}
		defer toimg.Close()

		jpeg.Encode(toimg, dst, &jpeg.Options{90})
	}
}

func createImages(id string, width string, height string) (path string) {
	fullSizePath := "public/" + id
	resizedPath := "public/" + id + "_" + width + "x" + height + ".jpg"

	downloadAndSaveOriginal(fullSizePath, id)
	createResizedImage(fullSizePath, resizedPath, width, height)
	return resizedPath
}

func widthHeightHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	width := params["width"]
	height := params["height"]

	resizedPath := createImages(id, width, height)
	http.ServeFile(w, r, resizedPath)
}

func squareHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	width := params["width"]
	height := params["width"]

	resizedPath := createImages(id, width, height)
	http.ServeFile(w, r, resizedPath)
}
