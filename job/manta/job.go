package mantajob

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/wanelo/image-server/job"
	"github.com/wanelo/image-server/uploader"
	"github.com/wanelo/image-server/uploader/manta/client"
)

type Job struct {
	Input          io.Reader
	Outputs        string
	BatchSize      int
	JobID          string
	InputsCount    int
	RemoteBasePath string
}

// CreateJob
func CreateJob(outputs string, remoteBasePath string, input io.Reader) (j *Job, err error) {
	j = &Job{Input: input, Outputs: outputs}

	mantaClient := client.DefaultClient()
	mantaJob := job.MantaJob{BasePath: remoteBasePath, Outputs: outputs}
	j.JobID, err = mantaJob.CreateMantaJob(mantaClient)
	if err != nil {
		return nil, err
	}

	j.InputsCount, err = j.createInputs(input)

	return j, err
}

fun (job Job) ToJobInput() io.Reader {

}

func (job Job) AddInputs() error {
	input := job.ToJobInput()

	err = mantaClient.AddJobInput(job.JobID, input)
	if err != nil {
		return err
	}

	err = mantaClient.EndJobInput(job.JobID)
	if err != nil {
		return err
	}
	return nil
}

func (job Job) AddJobs(partitionCount int) error {
	uploader := uploader.DefaultUploader(job.RemoteBasePath)
	remoteDirectory := fmt.Sprintf("stor/images/batches/%s", job.JobID)
	err := uploader.CreateDirectory(remoteDirectory)
	if err != nil {
		log.Println("Can't create remote directory", remoteDirectory, err)
		return err
	}

	for i := 0; i <= partitionCount; i++ {
		uploadBatchPartition(job.JobID, i, uploader)
	}

	mantaClient := client.DefaultClient()
	partitionPaths := ""
	for i := 0; i <= partitionCount; i++ {
		partitionPaths += fmt.Sprintf("/wanelo/stor/images/batches/%s/%d.txt\n", job.JobID, i)
	}
	log.Println(partitionPaths)

	err = mantaClient.AddJobInput(job.JobID, strings.NewReader(partitionPaths))
	if err != nil {
		return err
	}

	err = mantaClient.EndJobInput(job.JobID)
	if err != nil {
		return err
	}
	return nil
}

// createInputs splits input into multiple partitioned files with batch size
// and returns the number of partitions created.
// files are stored on tmp directory under job uuid
func (job Job) createInputs(input io.Reader) (partition int, err error) {
	dirName := fmt.Sprintf("tmp/%s", job.JobID)
	os.MkdirAll(dirName, 0700)

	reader := bufio.NewReader(input)
	count := 0
	partition = 0
	eof := false
	var items []string

	for !eof {
		var line string
		line, err = reader.ReadString('\n')
		if err == io.EOF {
			err = nil
			eof = true
		} else if err != nil {
			return partition, err
		}

		items = append(items, line)

		count++
		if count >= job.BatchSize || eof {
			count = 0
			err = writeBatchPartition(job.JobID, partition, items)
			if err != nil {
				log.Println("Can't write batch partition", job.JobID, partition, err)
				return 0, err
			}
			items = nil
			partition++
		}
	}

	return partition, nil
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
