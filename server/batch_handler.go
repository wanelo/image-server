package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"code.google.com/p/go-uuid/uuid"

	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/uploader"
)

func BatchHandler(w http.ResponseWriter, req *http.Request, sc *core.ServerConfiguration) {
	name := uuid.NewRandom()
	dirName := fmt.Sprintf("tmp/%s", name)
	os.MkdirAll(dirName, 0700)
	reader := bufio.NewReader(req.Body)
	uploader := uploader.DefaultUploader(sc.RemoteBasePath)

	var err error
	count := 0
	partition := 0
	eof := false
	var items []string

	for !eof {
		var line string
		line, err = reader.ReadString('\n')
		if err == io.EOF {
			err = nil
			eof = true
		} else if err != nil {
			return
		}

		items[count] = line

		count++
		if count >= 1000 {
			count = 0
			path := fmt.Sprintf("%s/%d.txt", dirName, partition)
			writeBatchPartition(path, items)
			uploadBatchPartition(path, uploader)
			partition++
		}
	}

}

func uploadBatchPartition(path string, uploader *uploader.Uploader) error {
	return nil
}

func writeBatchPartition(path string, lines []string) error {
	file, err := os.Create(path)

	if err != nil {
		return err
	}
	defer file.Close()

	for _, item := range lines {
		_, err = file.WriteString(strings.TrimSpace(item) + "\n")
		if err != nil {
			log.Println(err)
			break
		}
	}
	return err
}
