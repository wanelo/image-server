package server_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	"github.com/wanelo/image-server/core"
	fetcher "github.com/wanelo/image-server/fetcher/http"
	"github.com/wanelo/image-server/logger"
	"github.com/wanelo/image-server/logger/graphite"
	"github.com/wanelo/image-server/paths"
	"github.com/wanelo/image-server/server"
	"github.com/wanelo/image-server/uploader"

	. "github.com/wanelo/image-server/test"
)

func TestNewImageHandlerWithS3(t *testing.T) {
	deleteS3TestDirectory()

	sc := buildTestServerConfiguration()
	uploader.Initialize(sc)

	router := server.NewRouter(sc)

	request, _ := http.NewRequest("POST", "/test_namespace?outputs=x300.jpg,x300.webp&source=http%3A%2F%2Fcdn-s3-3.wanelo.com%2Fproduct%2Fimage%2F15209365%2Fx354.jpg", nil)
	response := httptest.NewRecorder()
	log.Println(sc)

	router.ServeHTTP(response, request)

	url := "https://s3.amazonaws.com/wanelo-dev/test/test_namespace/6da/b5f/6d8/d4bddc73fdff34d4f0507f7/info.json"
	resp, err := http.Head(url)
	Ok(t, err)
	Equals(t, "200 OK", resp.Status)
	Equals(t, "application/json", resp.Header.Get("Content-Type"))

	url = "https://s3.amazonaws.com/wanelo-dev/test/test_namespace/6da/b5f/6d8/d4bddc73fdff34d4f0507f7/x300.jpg"
	resp, err = http.Head(url)
	Ok(t, err)
	Equals(t, "200 OK", resp.Status)
	Equals(t, "image/jpeg", resp.Header.Get("Content-Type"))

	url = "https://s3.amazonaws.com/wanelo-dev/test/test_namespace/6da/b5f/6d8/d4bddc73fdff34d4f0507f7/x300.webp"
	resp, err = http.Head(url)
	Ok(t, err)
	Equals(t, "200 OK", resp.Status)
	Equals(t, "image/webp", resp.Header.Get("Content-Type"))
}

func deleteS3TestDirectory() {
	auth := aws.Auth{
		AccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey: os.Getenv("AWS_SECRET_KEY"),
	}
	client := s3.New(auth, aws.USEast)
	bucket := client.Bucket(os.Getenv("AWS_BUCKET"))
	bucket.Del("/test/test_namespace/6da/b5f/6d8/d4bddc73fdff34d4f0507f7/info.json")
	bucket.Del("/test/test_namespace/6da/b5f/6d8/d4bddc73fdff34d4f0507f7/original")
	bucket.Del("/test/test_namespace/6da/b5f/6d8/d4bddc73fdff34d4f0507f7/x300.jpg")
	bucket.Del("/test/test_namespace/6da/b5f/6d8/d4bddc73fdff34d4f0507f7/x300.webp")
}

func buildTestServerConfiguration() *core.ServerConfiguration {
	sc := &core.ServerConfiguration{
		LocalBasePath:  "../public",
		RemoteBasePath: "test",
		DefaultQuality: 90,
		AWSAccessKeyID: os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretKey:   os.Getenv("AWS_SECRET_KEY"),
		AWSBucket:      os.Getenv("AWS_BUCKET"),
	}

	loggers := []core.Logger{
		graphite.New(sc.GraphiteHost, sc.GraphitePort),
	}

	adapters := &core.Adapters{
		Fetcher: &fetcher.Fetcher{},
		Paths:   &paths.Paths{LocalBasePath: sc.LocalBasePath, RemoteBasePath: sc.RemoteBasePath, RemoteBaseURL: sc.RemoteBaseURL},
		Logger:  &logger.Logger{Loggers: loggers},
	}
	sc.Adapters = adapters
	return sc
}
