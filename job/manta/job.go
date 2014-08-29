package mantajob

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/wanelo/image-server/uploader"
	"github.com/wanelo/image-server/uploader/manta/client"
)

type Job struct {
	Outputs     string
	BatchSize   int
	JobID       string
	InputsCount int
}

func CreateJob(outputs string, batchSize int, input io.Reader) (job *Job, err error) {
	job = &Job{
		Outputs:   outputs,
		BatchSize: batchSize,
	}
	job.JobID, err = createBatchJob(outputs)
	if err != nil {
		return nil, err
	}

	job.InputsCount, err = job.createInputs(input)

	return job, err
}

func (job Job) AddJobs(partitionCount int, remoteBasePath string) error {
	uploader := uploader.DefaultUploader(remoteBasePath)
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
		if count >= job.BatchSize {
			count = 0
			writeBatchPartition(job.JobID, partition, items)
			items = nil
			partition++
		}
	}

	// write the remaining items
	if items != nil {
		err = writeBatchPartition(job.JobID, partition, items)
		if err != nil {
			log.Println("Can't write batch partition", job.JobID, partition, err)
			return partition, err
		}
	}

	return partition, nil
}

func createBatchJob(outputs string) (string, error) {
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

	opts := client.CreateJobOpts{Phases: phases}
	jobID, err := mantaClient.CreateJob(opts)
	if err != nil {
		return "", err
	}
	return jobID, nil
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
