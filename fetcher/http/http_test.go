package http_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	httpFetcher "github.com/image-server/image-server/fetcher/http"

	. "github.com/image-server/image-server/test"
)

func TestUniqueFetcher(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		fmt.Fprintln(w, `there is some content`)
	}))
	defer ts.Close()

	f := &httpFetcher.Fetcher{}

	defer os.Remove("valid")
	err := f.Fetch(ts.URL, "valid")

	Ok(t, err)
}

func TestUniqueFetcherOnEmptyFiles(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		fmt.Fprintf(w, ``)
	}))
	defer ts.Close()

	f := &httpFetcher.Fetcher{}

	defer os.Remove("blank.jpg")
	err := f.Fetch(ts.URL, "blank.jpg")

	Equals(t, "File is empty", fmt.Sprintf("%s", err))
}

func TestURLEscaping(t *testing.T) {
	path := "//hell[o]/(x)//two%20words/boo.jpg?something=fo(o)"
	expectedPath := "/hell[o]/(x)//two%20words/boo.jpg?something=fo(o)"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != expectedPath {
			t.Fail()
		}
		w.Header().Set("Content-Type", "image/jpeg")
		fmt.Fprintln(w, `there is some content`)
	}))
	defer ts.Close()

	f := &httpFetcher.Fetcher{}
	defer os.Remove("valid")
	err := f.Fetch(ts.URL+path, "valid")

	Ok(t, err)
}
