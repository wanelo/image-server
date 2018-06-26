package file_garbage_collector

import (
	"log"
	"os"
	"time"

	"github.com/image-server/image-server/core"
	"path/filepath"
)

func Start(sc *core.ServerConfiguration) {
	go func() {
		absolutePath, err := filepath.Abs(sc.LocalBasePath)
		if err != nil {
			log.Printf("Error Starting File Cleaner - Unable to create absolute path [%s]", sc.LocalBasePath)
		} else {
			var stat, _ = os.Stat(absolutePath)
			if absolutePath != "" && stat.IsDir() {
				log.Printf("Starting File Cleaner on path [%s]", absolutePath)
				filepath.Walk(absolutePath, func(path string, info os.FileInfo, err error) error {
					if !info.IsDir() {
						age := time.Now().Sub(info.ModTime())
						log.Printf("Workspace Contains at Startup file [%s] size [%d] modTime [%s] age [%s]\n", path, info.Size(), info.ModTime(), age)
					}
					return nil
				})
				for range sc.CleanUpTicker.C {
					tickTime := time.Now()
					log.Printf("Starting [tickID: %v]\n", tickTime)
					filepath.Walk(absolutePath, func(path string, info os.FileInfo, err error) error {
						if !info.IsDir() {
							age := tickTime.Sub(info.ModTime())
							if age > sc.MaxFileAge {
								log.Printf("[tickID: %v] Deleting file [%s] size [%d] modTime [%s] age [%s]\n", tickTime, path, info.Size(), info.ModTime(), age)
								var err = os.Remove(path)
								if err != nil {
									log.Printf("[tickID: %v] Error deleting file [%s]\n", tickTime, path)
								}
							}
						}
						return nil
					})
					log.Printf("Finished [tickID: %v]\n", tickTime)
				}
			} else {
				log.Printf("Error Starting File Cleaner - Invalid walk path [%s]", absolutePath)
			}
		}
	}()
}