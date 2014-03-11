package main

import (
	/*	"code.google.com/p/graphics-go/graphics"*/
	/*	"fmt"*/
	"github.com/gorilla/mux"
	"github.com/rainycape/magick"
	/*	"image"*/
	/*	"image/jpeg"*/
	"io"
	"log"
	"net/http"
	"os"
	/*	"runtime"*/
	"strconv"
	"time"
)

func main() {
	/*	runtime.GOMAXPROCS(100)*/
	/*	fmt.Printf("GOMAXPROCS is %d\n", runtime.GOMAXPROCS(0))*/

	r := mux.NewRouter()
	r.HandleFunc("/product/image/{id:[0-9]+}/{width:[0-9]+}x{height:[0-9]+}.{format}", rectangleHandler).Methods("GET")
	r.HandleFunc("/product/image/{id:[0-9]+}/x{width:[0-9]+}.{format}", squareHandler).Methods("GET")
	r.HandleFunc("/product/image/{id:[0-9]+}/w{width:[0-9]+}.{format}", widthHandler).Methods("GET")
	r.HandleFunc("/product/image/{id:[0-9]+}/full_size.{format}", fullSizeHandler).Methods("GET")
	http.Handle("/", r)
	log.Println("Listening on port 7000...")
	http.ListenAndServe(":7000", nil)
}

func downloadAndSaveOriginal(path string, productId string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		start := time.Now()
		resp, err := http.Get("http://cdn-s3-2.wanelo.com/product/image/" + productId + "/original.jpg")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		out, err := os.Create(path)
		defer out.Close()

		/*	imgBody := resp.Body*/
		io.Copy(out, resp.Body)
		elapsed := time.Since(start)
		log.Printf("Took %s to download image: %s", elapsed, path)
	}
}

/*func createJPG(fullSizePath string, resizedPath string, width string, height string) {
	file, err := os.Open(fullSizePath)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

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

	jpgOptions := jpeg.Options{95}
	jpeg.Encode(toimg, dst, &jpgOptions)
}*/

func createWithMagick(fullSizePath string, resizedPath string, width string, height string, format string) {
	start := time.Now()
	im, err := magick.DecodeFile(fullSizePath)
	if err != nil {
		log.Panicln(err)
		return
	}
	defer im.Dispose()

	w, _ := strconv.Atoi(width)
	h, _ := strconv.Atoi(height)

	im2, err := im.CropResize(w, h, magick.FHamming, magick.CSCenter)
	if err != nil {
		log.Panicln(err)
		return
	}

	out, err := os.Create(resizedPath)
	defer out.Close()

	info := magick.NewInfo()
	info.SetQuality(75)
	info.SetFormat(format)
	err = im2.Encode(out, info)

	if err != nil {
		log.Panicln(err)
		return
	}
	elapsed := time.Since(start)
	log.Printf("Took %s to generate image: %s", elapsed, resizedPath)
}

func createImages(id string, width string, height string, format string) (path string) {
	fullSizePath := "public/" + id
	resizedPath := "public/" + id + "_" + width + "x" + height + "." + format

	if _, err := os.Stat(resizedPath); os.IsNotExist(err) {
		downloadAndSaveOriginal(fullSizePath, id)
		/*	if format == "jpg" {*/
		/*		createJPG(fullSizePath, resizedPath, width, height)*/
		/*	} else if format == "webp" {*/
		createWithMagick(fullSizePath, resizedPath, width, height, format)
		/*	}*/
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

	if _, err := os.Stat(resizedPath); os.IsNotExist(err) {
		downloadAndSaveOriginal(fullSizePath, id)

		/* here*/
		im, err := magick.DecodeFile(fullSizePath)
		if err != nil {
			log.Panicln(err)
			return
		}
		defer im.Dispose()

		widthInt, _ := strconv.Atoi(width)
		heightInt := 0

		im2, err := im.CropResize(widthInt, heightInt, magick.FHamming, magick.CSCenter)
		if err != nil {
			log.Panicln(err)
			return
		}

		out, err := os.Create(resizedPath)
		defer out.Close()

		info := magick.NewInfo()
		info.SetQuality(75)
		info.SetFormat(format)
		err = im2.Encode(out, info)

		if err != nil {
			log.Panicln(err)
			return
		}
	}
	http.ServeFile(w, r, resizedPath)
}

func fullSizeHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	format := params["format"]

	fullSizePath := "public/" + id
	resizedPath := "public/" + id + "_full_size." + format

	if _, err := os.Stat(resizedPath); os.IsNotExist(err) {
		downloadAndSaveOriginal(fullSizePath, id)

		im, err := magick.DecodeFile(fullSizePath)
		if err != nil {
			log.Panicln(err)
			return
		}
		defer im.Dispose()

		out, err := os.Create(resizedPath)
		defer out.Close()

		info := magick.NewInfo()
		info.SetQuality(75)
		info.SetFormat(format)
		err = im.Encode(out, info)

		if err != nil {
			log.Panicln(err)
			return
		}
	}

	http.ServeFile(w, r, resizedPath)
}
