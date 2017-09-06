package server_test

import (
	"bytes"
	"github.com/image-server/image-server/core"
	fetcher "github.com/image-server/image-server/fetcher/http"
	"github.com/image-server/image-server/paths"
	"github.com/image-server/image-server/server"
	"github.com/image-server/image-server/uploader"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	. "github.com/image-server/image-server/test"
)

func TestNewFileHandlerWithData(t *testing.T) {
	sc := buildTestS3ServerConfiguration2()
	uploader.Initialize(sc)

	router := server.NewRouter(sc)
	uri := "/test_namespace/31e/8b3/187/a9f63f26d58c88bf09a7bbd/test.txt"
	var params map[string]string
	request, err := newFileUploadRequest(uri, params, "source", "../test/process.txt")
	Ok(t, err)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	json := ReaderToString(response.Body)

	Matches(t, "{}", json)
	ExpectFile(t, "../public/test_namespace/31e/8b3/187/a9f63f26d58c88bf09a7bbd/test.txt")
	log.Println(response.Body)
}

func buildTestS3ServerConfiguration2() *core.ServerConfiguration {
	sc := &core.ServerConfiguration{
		LocalBasePath:  "../public",
		RemoteBasePath: "test",
		DefaultQuality: 90,
	}

	adapters := &core.Adapters{
		Fetcher: &fetcher.Fetcher{},
		Paths:   &paths.Paths{LocalBasePath: sc.LocalBasePath, RemoteBasePath: sc.RemoteBasePath, RemoteBaseURL: sc.RemoteBaseURL},
	}
	sc.Adapters = adapters
	return sc
}

func newFileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return http.NewRequest("POST", uri, body)
}
