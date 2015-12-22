package mantajob

import (
	"fmt"
	"io"

	"github.com/image-server/image-server/job"
	"github.com/image-server/image-server/uploader/manta/client"
)

type Job struct {
	Input          io.Reader
	Outputs        string
	BatchSize      int
	JobID          string
	InputsCount    int
	RemoteBasePath string
	Namespace      string
	JobInput       io.Reader
}

// CreateJob takes the supplied Job metadata and image hash stream and initializes
// a Manta job. It then adds the inputs to the Manta job.
func CreateJob(outputs string, remoteBasePath string, namespace string, input io.Reader) (j *Job, err error) {
	mantaClient := client.DefaultClient()
	basePath := fmt.Sprintf("/%s/%s", mantaClient.User, remoteBasePath)
	j = &Job{Input: input, Outputs: outputs, RemoteBasePath: basePath, Namespace: namespace}

	mantaJob := job.MantaJob{BasePath: j.RemoteBasePath, Outputs: outputs, Namespace: namespace}
	j.JobID, err = mantaJob.CreateMantaJob(mantaClient)
	if err != nil {
		return nil, err
	}
	j.JobInput = j.ToMantaJobInput()

	return j, nil
}

func (j Job) AddInputs() error {
	mantaClient := client.DefaultClient()
	err := mantaClient.AddJobInput(j.JobID, j.JobInput)
	if err != nil {
		return err
	}

	err = mantaClient.EndJobInput(j.JobID)
	if err != nil {
		return err
	}
	return nil
}

// ToJobInput takes the input uploaded into the server and
// uses HashesToPaths to convert image hashes into manta paths
func (j Job) ToMantaJobInput() io.Reader {
	return job.HashesToPaths(j.Input, j.RemoteBasePath, j.Namespace)
}
