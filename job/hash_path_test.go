package job_test

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/wanelo/image-server/job"
	. "github.com/wanelo/image-server/test"
)

func TestToMantaPath(t *testing.T) {
	path := job.HashPath{BasePath: "/var", Namespace: "p", Hash: "6ad5544baa6f5e852e1af26f8c2e45db"}
	Equals(t, "/var/p/6ad/554/4ba/a6f5e852e1af26f8c2e45db", path.ToMantaPath())
}

func TestJoin(t *testing.T) {
	path := job.HashPath{BasePath: "/var", Namespace: "p", Hash: "6ad5544baa6f5e852e1af26f8c2e45db"}
	Equals(t, "/var/p/6ad/554/4ba/a6f5e852e1af26f8c2e45db/info.json", path.Join("info.json").ToMantaPath())
}

func TestMultiJoin(t *testing.T) {
	path := job.HashPath{BasePath: "/var", Namespace: "p", Hash: "6ad5544baa6f5e852e1af26f8c2e45db"}
	Equals(t, "/var/p/6ad/554/4ba/a6f5e852e1af26f8c2e45db/gad/zooks", path.Join("gad").Join("zooks").ToMantaPath())
}

func TestBuilderStyleJoin(t *testing.T) {
	path := job.HashPath{BasePath: "/var", Namespace: "p", Hash: "6ad5544baa6f5e852e1af26f8c2e45db"}
	path.Join("tacos.jpg")
	Equals(t, "/var/p/6ad/554/4ba/a6f5e852e1af26f8c2e45db", path.ToMantaPath())
}

func TestFromReader(t *testing.T) {
	input := strings.NewReader("6ad5544baa6f5e852e1af26f8c2e45db\n11111111111111111111111111111111\n")
	hashPaths := job.HashPaths(input, "/var", "p")
	Equals(t, job.HashPath{BasePath: "/var", Namespace: "p", Hash: "6ad5544baa6f5e852e1af26f8c2e45db"}, hashPaths[0])
}

func TestHashPathJoin(t *testing.T) {
	path1 := job.HashPath{BasePath: "/var", Namespace: "p", Hash: "6ad5544baa6f5e852e1af26f8c2e45db"}
	path2 := job.HashPath{BasePath: "/home", Namespace: "w", Hash: "6ad5544baa6f5e852e1af26f8c2e45db"}
	list := []job.HashPath{path1, path2}
	output := job.HashPathJoin(list, "original")
	Equals(t, "/var/p/6ad/554/4ba/a6f5e852e1af26f8c2e45db/original", output[0].ToMantaPath())
}

func TestHashPathReader(t *testing.T) {
	path1 := job.HashPath{BasePath: "/var", Namespace: "p", Hash: "6ad5544baa6f5e852e1af26f8c2e45db"}
	path2 := job.HashPath{BasePath: "/home", Namespace: "w", Hash: "6ad5544baa6f5e852e1af26f8c2e45db"}
	list := []job.HashPath{path1, path2}
	list = job.HashPathJoin(list, "original")
	output := job.HashPathReader(list)
	reader := bufio.NewReader(output)
	line, _ := reader.ReadString('\n')
	Equals(t, "/var/p/6ad/554/4ba/a6f5e852e1af26f8c2e45db/original\n", line)
}

func TestHashesToPaths(t *testing.T) {
	input := strings.NewReader("6ad5544baa6f5e852e1af26f8c2e89bc\n11111111111111111111111111111111\n")
	expected := "/var/p/6ad/554/4ba/a6f5e852e1af26f8c2e89bc/original\n/var/p/111/111/111/11111111111111111111111/original\n"

	paths := job.HashesToPaths(input, "/var", "p")
	Equals(t, expected, readerToString(paths))
}

func readerToString(input io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(input)
	return buf.String()
}
