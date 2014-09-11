package server

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/wanelo/image-server/core"
	mantajob "github.com/wanelo/image-server/job/manta"
	"github.com/wanelo/image-server/uploader/manta/client"
)

func CreateBatchHandler(w http.ResponseWriter, req *http.Request, sc *core.ServerConfiguration) {
	vars := mux.Vars(req)
	namespace := vars["namespace"]

	r := render.New(render.Options{
		IndentJSON: true,
	})

	job, err := mantajob.CreateJob(sc.Outputs, sc.RemoteBasePath, namespace, req.Body)
	if err != nil {
		errorHandlerJSON(err, w, r, http.StatusInternalServerError)
		return
	}

	json := map[string]string{
		"job_id": job.JobID,
	}

	r.JSON(w, http.StatusOK, json)
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
		result, err := getJobOutput(uuid, mantaClient)
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

func getJobOutput(uuid string, mantaClient *client.Client) (io.Reader, error) {
	output, err := mantaClient.GetJobOutput(uuid)
	if err != nil {
		return nil, err
	}

	result, err := mantaClient.GetObject(output)
	if err != nil {
		return nil, err
	}

	return result, nil
}
