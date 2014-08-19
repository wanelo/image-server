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

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/uploader"
	"github.com/wanelo/image-server/uploader/manta/client"
)

func CreateBatchHandler(w http.ResponseWriter, req *http.Request, sc *core.ServerConfiguration) {
	r := render.New(render.Options{
		IndentJSON: true,
	})

	batchSize := 100
	name := uuid.NewRandom().String()
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
			errorHandlerJSON(err, w, r, http.StatusInternalServerError)
			return
		}

		items = append(items, line)

		count++
		if count >= batchSize {
			count = 0
			writeBatchPartition(name, partition, items)
			items = nil
			partition++
		}
	}

	// write the remaining items
	if items != nil {
		err = writeBatchPartition(name, partition, items)
		if err != nil {
			log.Println("Can't write batch partition", name, partition, err)
			errorHandlerJSON(err, w, r, http.StatusInternalServerError)
			return
		}
	}

	remoteBasePath := fmt.Sprintf("stor/images/batches/%s", name)
	err = uploader.CreateDirectory(remoteBasePath)
	if err != nil {
		log.Println("Can't create remote directory", remoteBasePath, err)
		errorHandlerJSON(err, w, r, http.StatusInternalServerError)
		return
	}

	for i := 0; i <= partition; i++ {
		uploadBatchPartition(name, i, uploader)
	}

	jobID, err := createBatchJob(name, partition)
	if err != nil {
		log.Println("Can't initialize manta job:", err)
		errorHandlerJSON(err, w, r, http.StatusInternalServerError)
		return
	}

	json := map[string]string{
		"job_id": jobID,
	}

	r.JSON(w, http.StatusOK, json)
}

func BatchHandler(w http.ResponseWriter, req *http.Request, sc *core.ServerConfiguration) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]

	mantaClient := client.DefaultClient()

	// if not complete return 202
	job, err := mantaClient.GetJob(uuid)
	if err != nil {
		errorHandler(err, w, req, 500)
	}

	if job.State == "done" {
		fmt.Fprint(w, "YAY done")
		w.WriteHeader(200)
	} else {
		w.WriteHeader(202)
	}
}

func uploadBatchPartition(jobID string, partition int, uploader *uploader.Uploader) error {
	localPath := fmt.Sprintf("tmp/%s/%d.txt", jobID, partition)
	remotePath := fmt.Sprintf("stor/images/batches/%s/%d.txt", jobID, partition)
	log.Printf("Uploading batch from %s to %s", localPath, remotePath)

	uploader.Upload(localPath, remotePath)
	return nil
}

func writeBatchPartition(jobID string, partition int, lines []string) error {
	path := fmt.Sprintf("tmp/%s/%d.txt", jobID, partition)
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

func createBatchJob(uuid string, partitionCount int) (string, error) {
	mantaClient := client.DefaultClient()

	phases := []client.Phase{
		{Type: "map",
			Exec: "/assets/wanelo/public/images/bin/images-solaris-1.0.6 --remote_base_path public/images --outputs full_size.jpg,full_size.webp,x110-q90.jpg,x200-q90.jpg,x354-q80.jpg,w620-q80.jpg,w736-q75.jpg,w1472-q65.jpg,x110-q90.webp,x200-q90.webp,x354-q80.webp,w620-q80.webp,w736-q75.webp,w1472-q65.webp process",
			Init: "/assets/wanelo/stor/images/init.sh",
			Assets: []string{
				"/wanelo/public/images/bin/images-solaris-1.0.6",
				"/wanelo/stor/images/init.sh",
			},
		},
		{Type: "reduce", Exec: "cat"},
	}
	opts := client.CreateJobOpts{Name: uuid, Phases: phases}
	jobID, err := mantaClient.CreateJob(opts)
	if err != nil {
		return "", err
	}

	partitionPaths := ""
	for i := 0; i <= partitionCount; i++ {
		partitionPaths += fmt.Sprintf("/wanelo/stor/images/batches/%s/%d.txt\n", uuid, i)
	}
	log.Println(partitionPaths)

	err = mantaClient.AddJobInput(jobID, strings.NewReader(partitionPaths))
	if err != nil {
		return "", err
	}

	err = mantaClient.EndJobInput(jobID)
	if err != nil {
		return "", err
	}

	return jobID, nil
}
