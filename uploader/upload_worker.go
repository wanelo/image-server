package uploader

type UploadWork struct {
	Source      string
	Destination string
	Func        func(string, string) error
}

func UploadWorkers(concurrency uint) chan *UploadWork {

	jobs := make(chan *UploadWork)

	for i := uint(0); i < concurrency; i++ {
		go uploadWorker(jobs)
	}

	return jobs
}

func uploadWorker(in chan *UploadWork) {
	for {
		t := <-in
		t.Func(t.Source, t.Destination)
	}
}
