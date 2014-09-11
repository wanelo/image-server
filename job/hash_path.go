package job

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type HashPath struct {
	BasePath  string
	Namespace string
	Hash      string
	parts     []string
}

func HashesToPaths(input io.Reader, basePath string, namespace string) io.Reader {
	hashPaths := HashPaths(input, basePath, namespace)
	list := HashPathJoin(hashPaths, "original")
	return HashPathReader(list)
}

func HashPaths(input io.Reader, basePath string, namespace string) []HashPath {
	reader := bufio.NewReader(input)
	hashPaths := []HashPath{}

	line, err := reader.ReadString('\n')
	for err != io.EOF {
		hashPaths = append(hashPaths, HashPath{
			BasePath:  basePath,
			Namespace: namespace,
			Hash:      strings.TrimSpace(line),
		})
		line, err = reader.ReadString('\n')
	}

	return hashPaths
}

func HashPathJoin(list []HashPath, joined string) []HashPath {
	output := []HashPath{}
	for _, item := range list {
		output = append(output, *item.Join(joined))
	}
	return output
}

func HashPathReader(list []HashPath) io.Reader {
	var buffer bytes.Buffer
	for _, item := range list {
		buffer.WriteString(item.ToMantaPath() + "\n")
	}
	return strings.NewReader(buffer.String())
}

func (h *HashPath) ToMantaPath() string {
	var partitions = []string{h.Hash[0:3], h.Hash[3:6], h.Hash[6:9], h.Hash[9:32]}
	var joinedPartitions = strings.Join(partitions, "/")
	var joinedParts = strings.Join(append([]string{""}, h.parts...), "/")
	return fmt.Sprintf("%s/%s/%s%s", h.BasePath, h.Namespace, joinedPartitions, joinedParts)
}

func (h *HashPath) Join(part string) *HashPath {
	return &HashPath{
		BasePath:  h.BasePath,
		Namespace: h.Namespace,
		Hash:      h.Hash,
		parts:     append(h.parts, part),
	}
}
