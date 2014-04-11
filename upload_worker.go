package main

type UploadWork struct {
	ic *ImageConfiguration
}

func UploadWorker(in chan *UploadWork, fn func(*ImageConfiguration)) {
	for {
		t := <-in
		fn(t.ic)
	}
}

func UploadWorkers(fn func(*ImageConfiguration), concurrency int) chan *UploadWork {
	jobs := make(chan *UploadWork)

	for i := 0; i < concurrency; i++ {
		go UploadWorker(jobs, fn)
	}

	return jobs
}
