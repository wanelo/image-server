package uploader

import "github.com/wanelo/image-server/core"

type UploadWork struct {
	ImageConfiguration *core.ImageConfiguration
	Func               func(*core.ImageConfiguration)
}

func UploadWorker(in chan *UploadWork) {
	for {
		t := <-in
		t.Func(t.ImageConfiguration)
	}
}

func UploadWorkers(concurrency uint) chan *UploadWork {

	jobs := make(chan *UploadWork)

	for i := uint(0); i < concurrency; i++ {
		go UploadWorker(jobs)
	}

	return jobs
}
