package uploader

import "github.com/wanelo/image-server/core"

type UploadWork struct {
	ImageConfiguration *core.ImageConfiguration
}

func UploadWorker(in chan *UploadWork, fn func(*core.ImageConfiguration)) {
	for {
		t := <-in
		fn(t.ImageConfiguration)
	}
}

func UploadWorkers(fn func(*core.ImageConfiguration), concurrency int) chan *UploadWork {
	jobs := make(chan *UploadWork)

	for i := 0; i < concurrency; i++ {
		go UploadWorker(jobs, fn)
	}

	return jobs
}
