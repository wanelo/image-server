package request_test

import (
	"testing"

	"github.com/image-server/image-server/request"
	. "github.com/image-server/image-server/test"
)

type FakeUploader struct{}

func (u FakeUploader) CreateDirectory(string) error           { return nil }
func (u FakeUploader) Upload(string, string, string) error    { return nil }
func (u FakeUploader) ListDirectory(string) ([]string, error) { return []string{"a", "b"}, nil }
func (u FakeUploader) Initialize() error                      { return nil }

type FakePaths struct{}

func (u FakePaths) LocalInfoPath(string, string) string                                  { return "" }
func (u FakePaths) RemoteInfoPath(string, string) string                                 { return "" }
func (u FakePaths) TempImagePath(string) string                                          { return "" }
func (u FakePaths) RandomTempPath() string                                               { return "" }
func (u FakePaths) LocalOriginalPath(string, string) string                              { return "" }
func (u FakePaths) LocalImageDirectory(string, string) string                            { return "" }
func (u FakePaths) LocalImagePath(namespace string, md5 string, imageName string) string { return "" }
func (u FakePaths) RemoteImageDirectory(namespace string, md5 string) string {
	return "images/p/f94/4de/077/34f1868a4355e1b86052704"
}
func (u FakePaths) RemoteOriginalPath(string, string) string                              { return "" }
func (u FakePaths) RemoteOriginalURL(string, string) string                               { return "" }
func (u FakePaths) RemoteImagePath(namespace string, md5 string, imageName string) string { return "" }

func sampleRequest() *request.Request {
	return &request.Request{
		Uploader:  &FakeUploader{},
		Paths:     &FakePaths{},
		Namespace: "p",
		Outputs:   []string{"a", "b", "c"},
		Hash:      "f944de07734f1868a4355e1b86052704",
	}
}

func TestCalculateMissingOutputs(t *testing.T) {
	r := sampleRequest()
	missing, err := r.CalculateMissingOutputs()
	Ok(t, err)
	Equals(t, []string{"c"}, missing)
}

func TestCalculateMissingOutputsWhenOutputIsNotPresent(t *testing.T) {
	r := sampleRequest()
	r.Outputs = []string{}

	missing, err := r.CalculateMissingOutputs()
	Ok(t, err)
	Equals(t, []string(nil), missing)
}
