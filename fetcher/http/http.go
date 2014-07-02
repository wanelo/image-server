package http

import (
	"fmt"
	"io"
	"log"
	gohttp "net/http"
	"os"
	"path/filepath"
	"time"
)

type Fetcher struct{}

func (f *Fetcher) Fetch(url string, destination string) error {
	if _, err := os.Stat(destination); os.IsNotExist(err) {
		start := time.Now()

		resp, err := gohttp.Get(url)

		if err != nil || resp.StatusCode != 200 {
			log.Printf("Unable to download image: %s, status code: %d", url, resp.StatusCode)
			log.Println(err)
			return fmt.Errorf("Unable to download image: %s, status code: %d", url, resp.StatusCode)
		}
		log.Printf("Downloaded from %s with code %d", url, resp.StatusCode)
		defer resp.Body.Close()

		dir := filepath.Dir(destination)
		os.MkdirAll(dir, 0700)

		out, err := os.Create(destination)
		defer out.Close()
		if err != nil {
			log.Printf("Unable to create file: %s", destination)
			log.Println(err)
			return fmt.Errorf("Unable to create file: %s", destination)
		}

		io.Copy(out, resp.Body)
		log.Printf("Took %s to download image: %s", time.Since(start), destination)
	}
	return nil
}
