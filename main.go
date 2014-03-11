package main

import (
	"code.google.com/p/graphics-go/graphics"
	"github.com/gorilla/mux"
	"github.com/quirkey/magick"
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
	r.HandleFunc("/product/image/{id:[0-9]+}/{width:[0-9]+}x{height:[0-9]+}.{format}", rectangleHandler).Methods("GET")
	r.HandleFunc("/product/image/{id:[0-9]+}/x{width:[0-9]+}.{format}", squareHandler).Methods("GET")
	r.HandleFunc("/product/image/{id:[0-9]+}/w{width:[0-9]+}.{format}", widthHandler).Methods("GET")
	r.HandleFunc("/product/image/{id:[0-9]+}/full_size.{format}", fullSizeHandler).Methods("GET")
	http.Handle("/", r)
	log.Println("Listening on port 3000...")
	http.ListenAndServe(":3000", nil)
}

func downloadAndSaveOriginal(path string, productId string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		resp, err := http.Get("http://cdn-s3-2.wanelo.com/product/image/" + productId + "/original.jpg")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		out, err := os.Create(path)
		defer out.Close()

		/*	imgBody := resp.Body*/
		io.Copy(out, resp.Body)
	}
}

func createJPG(fullSizePath string, resizedPath string, width string, height string) {
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

func createWEBP(fullSizePath string, resizedPath string, width string, height string) {
	image, err := magick.NewFromFile(fullSizePath)

	err = image.Resize(height + "x" + width + "^")
	if err != nil {
		log.Fatal(err)
	}
	err = image.SetProperty("gravity", "center")
	if err != nil {
		log.Fatal(err)
	}
	err = image.SetProperty("extent", height+"x"+width)
	if err != nil {
		log.Fatal(err)
	}
	/*	image.Quality(85)*/
	/*	err = image.Crop(width + "x" + height + "")*/

	if err != nil {
		log.Fatal(err)
	}
	err = image.ToFile(resizedPath)
	if err != nil {
		log.Fatal(err)
	}
}

func createImages(id string, width string, height string, format string) (path string) {
	fullSizePath := "public/" + id
	resizedPath := "public/" + id + "_" + width + "x" + height + "." + format

	downloadAndSaveOriginal(fullSizePath, id)
	if format == "jpg" {
		createJPG(fullSizePath, resizedPath, width, height)
	} else if format == "webp" {
		createWEBP(fullSizePath, resizedPath, width, height)
	}

	return resizedPath
}

func rectangleHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	width := params["width"]
	height := params["height"]
	format := params["format"]

	resizedPath := createImages(id, width, height, format)
	http.ServeFile(w, r, resizedPath)
}

func squareHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	width := params["width"]
	height := params["width"]
	format := params["format"]

	resizedPath := createImages(id, width, height, format)
	http.ServeFile(w, r, resizedPath)
}

func widthHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	width := params["width"]
	format := params["format"]

	fullSizePath := "public/" + id
	resizedPath := "public/" + id + "_w" + width + "." + format

	downloadAndSaveOriginal(fullSizePath, id)

	image, err := magick.NewFromFile(fullSizePath)
	if err != nil {
		log.Fatal(err)
	}
	err = image.Resize(width)
	if err != nil {
		log.Fatal(err)
	}
	err = image.ToFile(resizedPath)
	if err != nil {
		log.Fatal(err)
	}
	http.ServeFile(w, r, resizedPath)
}

func fullSizeHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	format := params["format"]

	fullSizePath := "public/" + id
	resizedPath := "public/" + id + "_full_size." + format

	downloadAndSaveOriginal(fullSizePath, id)

	image, err := magick.NewFromFile(fullSizePath)
	if err != nil {
		log.Fatal(err)
	}
	err = image.ToFile(resizedPath)
	if err != nil {
		log.Fatal(err)
	}
	http.ServeFile(w, r, fullSizePath)
}
