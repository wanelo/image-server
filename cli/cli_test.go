package cli

import (
	"testing"

	. "github.com/wanelo/image-server/test"
)

func TestLineToItem(t *testing.T) {
	item, err := lineToItem("6ad5544baa6f5e852e1af26f8c2e45db http://example.com/image.url")
	Ok(t, err)
	Equals(t, "6ad5544baa6f5e852e1af26f8c2e45db", item.Hash)
	Equals(t, "http://example.com/image.url", item.URL)
}

func TestItemToTabDelimited(t *testing.T) {
	item := Item{"6ad5544baa6f5e852e1af26f8c2e45db", "http://example.com/image.jpg", 40, 30}
	expected := "6ad5544baa6f5e852e1af26f8c2e45db\thttp://example.com/image.jpg\t40\t30\n"
	Equals(t, expected, item.ToTabDelimited())
}
