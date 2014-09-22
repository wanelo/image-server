package cli_test

import (
	"testing"

	"github.com/wanelo/image-server/cli"
	. "github.com/wanelo/image-server/test"
)

func TestItemToHash(t *testing.T) {
	image := &cli.Image{LocalOriginalPath: "public/p/6ad/554/4ba/a6f5e852e1af26f8c2e45db/original"}

	Equals(t, "6ad5544baa6f5e852e1af26f8c2e45db", image.ToHash())
}
