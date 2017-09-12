package server_test

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/image-server/image-server/core"
	fetcher "github.com/image-server/image-server/fetcher/http"
	"github.com/image-server/image-server/paths"
	"github.com/image-server/image-server/server"
	"github.com/image-server/image-server/uploader"

	. "github.com/image-server/image-server/test"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"fmt"
)

func TestNewImageHandlerWithData(t *testing.T) {
	if !hasAwsAuthentication() {
		return
	}

	deleteS3TestDirectory()
	sc := buildTestS3ServerConfiguration()
	uploader.Initialize(sc)

	router := server.NewRouter(sc)
	uri := "/test_namespace?outputs=x300.jpg,x300.webp"
	imagePath := "../test/images/a.jpg"
	request, err := newUploadRequest(uri, imagePath)
	Ok(t, err)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	json := ReaderToString(response.Body)

	Matches(t, "\"hash\": \"31e8b3187a9f63f26d58c88bf09a7bbd\"", json)
	Matches(t, "\"height\": 496", json)
	Matches(t, "\"width\": 574", json)
	Matches(t, "\"content_type\": \"image/jpeg\"", json)
	log.Println(response.Body)
	deleteS3TestDirectory()
}

func TestNewImageHandlerWithS3(t *testing.T) {
	if !hasAwsAuthentication() {
		return
	}

	deleteS3TestDirectory()
	sc := buildTestS3ServerConfiguration()
	uploader.Initialize(sc)

	router := server.NewRouter(sc)

	request, _ := http.NewRequest("POST", "/test_namespace?outputs=x300.jpg,x300.webp&source=http%3A%2F%2Fcdn-s3-3.wanelo.com%2Fproduct%2Fimage%2F15209365%2Fx354.jpg", nil)
	response := httptest.NewRecorder()
	log.Println(sc)

	router.ServeHTTP(response, request)

	url := s3UrlForPath("test/test_namespace/6da/b5f/6d8/d4bddc73fdff34d4f0507f7/info.json")
	resp, err := http.Head(url)
	Ok(t, err)
	Equals(t, "200 OK", resp.Status)
	Equals(t, "application/json", resp.Header.Get("Content-Type"))

	url = s3UrlForPath("test/test_namespace/6da/b5f/6d8/d4bddc73fdff34d4f0507f7/x300.jpg")
	log.Println(url)
	resp, err = http.Head(url)
	Ok(t, err)
	Equals(t, "200 OK", resp.Status)
	Equals(t, "image/jpeg", resp.Header.Get("Content-Type"))

	url = s3UrlForPath("test/test_namespace/6da/b5f/6d8/d4bddc73fdff34d4f0507f7/x300.webp")
	resp, err = http.Head(url)
	Ok(t, err)
	Equals(t, "200 OK", resp.Status)
	Equals(t, "image/webp", resp.Header.Get("Content-Type"))
	deleteS3TestDirectory()
}
func s3UrlForPath(path string) string {
	return fmt.Sprintf("https://s3-%s.amazonaws.com/%s/%s", os.Getenv("AWS_REGION"), os.Getenv("AWS_BUCKET"), path)
}

func hasAwsAuthentication() bool {
	hasRegion := len(os.Getenv("AWS_REGION")) > 0
	hasBucket := len(os.Getenv("AWS_BUCKET")) > 0
	return hasRegion && hasBucket
}

func deleteS3TestDirectory() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}))
	svc := s3.New(sess)

	resp, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket:    aws.String(os.Getenv("AWS_BUCKET")),
		Prefix:    aws.String("test"),
	})
	if err == nil {
		entries := resp.Contents
		for _, entry := range entries {
			key := aws.StringValue(entry.Key)
			fmt.Printf("Deleting key: [%s]\n", key)
			svc.DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String(os.Getenv("AWS_BUCKET")),
				Key:    aws.String(key),
			})
		}
	}
}

func buildTestS3ServerConfiguration() *core.ServerConfiguration {
	sc := &core.ServerConfiguration{
		LocalBasePath:  "../public",
		RemoteBasePath: "test",
		DefaultQuality: 90,
		UploaderType: "aws",
		AWSBucket:      os.Getenv("AWS_BUCKET"),
		AWSRegion:      os.Getenv("AWS_REGION"),
	}

	adapters := &core.Adapters{
		Fetcher: &fetcher.Fetcher{},
		Paths:   &paths.Paths{LocalBasePath: sc.LocalBasePath, RemoteBasePath: sc.RemoteBasePath, RemoteBaseURL: sc.RemoteBaseURL},
	}
	sc.Adapters = adapters
	return sc
}

func buildTestServerConfiguration() *core.ServerConfiguration {
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

func newUploadRequest(uri string, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return http.NewRequest("POST", uri, file)
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
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
