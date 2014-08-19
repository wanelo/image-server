package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// Phase represents a task to be executed as part of a Job
type Phase struct {
	Type   string   `json:"type,omitempty"`   // Task type, one of 'map' or 'reduce' (optional)
	Assets []string `json:"assets,omitempty"` // An array of objects to be placed in the compute zones (optional)
	Exec   string   `json:"exec"`             // The actual shell statement to execute
	Init   string   `json:"init"`             // Shell statement to execute in each compute zone before any tasks are executed
	Count  int      `json:"count,omitempty"`  // If type is 'reduce', an optional number of reducers for this phase (default is 1)
	Memory int      `json:"memory,omitempty"` // Amount of DRAM to give to your compute zone (in Mb, optional)
	Disk   int      `json:"disk,omitempty"`   // Amount of disk space to give to your compute zone (in Gb, optional)
}

// CreateJobOpts represent the option that can be specified
// when creating a job.
type CreateJobOpts struct {
	Name   string  `json:"name,omitempty"` // Job Name (optional)
	Phases []Phase `json:"phases"`         // Tasks to execute as part of this job
}

// Job represents the status of a job.
type Job struct {
	Id                 string      // Job unique identifier
	Name               string      `json:"name,omitempty"` // Job Name
	State              string      // Job state
	Cancelled          bool        // Whether the job has been cancelled or not
	InputDone          bool        // Whether the inputs for the job is still open or not
	Stats              JobStats    `json:"stats,omitempty"` // Job statistics
	TimeCreated        string      // Time the job was created at
	TimeDone           string      `json:"timeDone,omitempty"`           // Time the job was completed
	TimeArchiveStarted string      `json:"timeArchiveStarted,omitempty"` // Time the job archiving started
	TimeArchiveDone    string      `json:"timeArchiveDone,omitempty"`    // Time the job archiving completed
	Phases             []Phase     `json:"phases"`                       // Job tasks
	Options            interface{} // Job options
}

// JobStats represents statistics about a job
type JobStats struct {
	Errors    int // Number or errors
	Outputs   int // Number of output produced
	Retries   int // Number of retries
	Tasks     int // Total number of task in the job
	TasksDone int // number of tasks done
}

// CreateJob submits a job to Manta and returns the URI for the job.
func (c *Client) CreateJob(opts CreateJobOpts) (string, error) {
	headers := make(http.Header)
	headers.Add("content-type", "application/json")
	json, err := json.Marshal(opts)
	if err != nil {
		return "", fmt.Errorf("Can't create json from opts: %v", err)
	}
	body := bytes.NewReader(json)

	resp, err := c.Post("jobs", headers, body)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	err = c.ensureStatus(resp, 201)
	if err != nil {
		return "", err
	}

	location := resp.Header.Get("Location")
	jobID := strings.Split(location, "/")[3]
	log.Println("Created Manta Job:", jobID)
	return jobID, nil
}

// AddJobInput submits input to a job previously created by CreateJob
func (c *Client) AddJobInput(jobID string, input io.Reader) error {
	headers := make(http.Header)
	headers.Add("content-type", "text/plain")

	path := fmt.Sprintf("jobs/%s/live/in", jobID)
	resp, err := c.Post(path, headers, input)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.ensureStatus(resp, 204)
}

// EndJobInput closes input for a job previously created by CreateJob
func (c *Client) EndJobInput(jobID string) error {
	path := fmt.Sprintf("jobs/%s/live/in/end", jobID)
	resp, err := c.Post(path, nil, nil)
	if err != nil {
		return err
	}

	return c.ensureStatus(resp, 202)
}

func (c *Client) GetJob(jobID string) (Job, error) {
	path := fmt.Sprintf("jobs/%s/live/status", jobID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return Job{}, err
	}

	err = c.ensureStatus(resp, 200)
	if err != nil {
		return Job{}, err
	}

	job := new(Job)
	err = json.NewDecoder(resp.Body).Decode(job)
	if err != nil {
		return Job{}, err
	}

	return *job, nil
}
