package main

import "testing"

func TestRemoteImageUrl(t *testing.T) {
	serverConfiguration, _ = loadServerConfiguration("config/test.json")
	ic := &ImageConfiguration{model: "product", imageType: "image", id: "55"}
	url := ic.RemoteImageUrl()
	expectedUrl := "http://cdn-s3-2.wanelo.com/product/image/55/original.jpg"

	if url != expectedUrl {
		t.Errorf("expected %s to be %s", url, expectedUrl)
	}
}
