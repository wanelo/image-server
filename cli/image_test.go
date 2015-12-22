package cli_test

import (
	"testing"

	"github.com/image-server/image-server/cli"
	. "github.com/image-server/image-server/test"
)

func TestItemToHash(t *testing.T) {
	image := &cli.Image{LocalOriginalPath: "public/p/6ad/554/4ba/a6f5e852e1af26f8c2e45db/original"}

	Equals(t, "6ad5544baa6f5e852e1af26f8c2e45db", image.ToHash())
}
