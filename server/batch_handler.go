package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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

	batchSize, err := strconv.Atoi(req.FormValue("batch_size"))
	if err != nil {
		batchSize = 1000
	}

	name := uuid.NewRandom().String()
	dirName := fmt.Sprintf("tmp/%s", name)
	os.MkdirAll(dirName, 0700)
	reader := bufio.NewReader(req.Body)
	uploader := uploader.DefaultUploader(sc.RemoteBasePath)

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

	jobID, err := createBatchJob(name, sc.Outputs)
	if err != nil {
		log.Println("Can't initialize manta job:", err)
		errorHandlerJSON(err, w, r, http.StatusInternalServerError)
		return
	}

	json := map[string]string{
		"job_id": jobID,
	}

	r.JSON(w, http.StatusOK, json)
	go addJobs(name, jobID, partition)
}

func BatchHandler(w http.ResponseWriter, req *http.Request, sc *core.ServerConfiguration) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]

	mantaClient := client.DefaultClient()
	job, err := mantaClient.GetJob(uuid)

	if err != nil {
		log.Println(err)
		errorHandler(err, w, req, 500)
		return
	}

	if job.State == "done" {

		var output string
		output, err = mantaClient.GetJobOutput(uuid)

		if err != nil {
			log.Println(err)
			errorHandler(err, w, req, 500)
			return
		}

		result, err := mantaClient.GetObject(output)
		if err != nil {
			log.Println(err)
			errorHandler(err, w, req, 500)
			return
		}

		w.WriteHeader(200)
		io.Copy(w, result)
	} else {
		// if not complete return 202
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

func createBatchJob(uuid string, outputs string) (string, error) {
	mantaClient := client.DefaultClient()
	remoteBasePath := "public/images"
	exec := fmt.Sprintf("/assets/wanelo/public/images/bin/images-solaris-1.0.6 --remote_base_path %s --outputs %s process", remoteBasePath, outputs)

	phases := []client.Phase{
		{Type: "map",
			Exec: exec,
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
	return jobID, nil
}

func addJobs(uuid string, jobID string, partitionCount int) error {
	mantaClient := client.DefaultClient()
	partitionPaths := ""
	for i := 0; i <= partitionCount; i++ {
		partitionPaths += fmt.Sprintf("/wanelo/stor/images/batches/%s/%d.txt\n", uuid, i)
	}
	log.Println(partitionPaths)

	err := mantaClient.AddJobInput(jobID, strings.NewReader(partitionPaths))
	if err != nil {
		return err
	}

	err = mantaClient.EndJobInput(jobID)
	if err != nil {
		return err
	}
	return nil
}
