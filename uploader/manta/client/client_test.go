package client_test

import (
	"testing"

	. "github.com/image-server/image-server/test"
	"github.com/image-server/image-server/uploader/manta/client"
)

func TestLineToItem(t *testing.T) {
	c := client.DefaultClient()

	if c.User == "" {
		// this test currently hits manta service, this test should only run on development machine
		return
	}

	entries, err := c.ListDirectory("/wanelo/public/images/p/000/018/d7f/5a0f2b3cfa4cdeb904c29c6/")
	Ok(t, err)

	entry := entries[0]
	Equals(t, entry.Name, "info.json")
	Equals(t, entry.Type, "object")

	// fmt.Println(entries)
}
