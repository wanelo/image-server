package core

import "testing"

func TestRemoteImageURL(t *testing.T) {
	sc := &ServerConfiguration{SourceDomain: "http://cdn-s3-2.wanelo.com"}
	ic := &ImageConfiguration{ServerConfiguration: sc, Model: "product", ImageType: "image", ID: "55"}
	url := ic.RemoteImageURL()
	expectedURL := "http://cdn-s3-2.wanelo.com/product/image/55/original.jpg"

	if url != expectedURL {
		t.Errorf("expected %s to be %s", url, expectedURL)
	}
}
