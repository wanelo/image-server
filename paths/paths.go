package paths

import (
	"crypto/md5"
	"fmt"
	"strings"
)

// Paths
type Paths struct {
	LocalBasePath  string
	RemoteBasePath string
	RemoteBaseURL  string
}

// LocalOriginalPath returns local path for original image
func (p *Paths) LocalOriginalPath(namespace string, md5 string) string {
	return fmt.Sprintf("%s/%s", p.LocalBasePath, p.originalPath(namespace, md5))
}

// LocalImageDirectory returns location for locally cached images
func (p *Paths) LocalImageDirectory(namespace string, md5 string) string {
	return fmt.Sprintf("%s/%s", p.LocalBasePath, p.imageDirectory(namespace, md5))
}

// RemoteImageDirectory returns location for directory for images and info
func (p *Paths) RemoteImageDirectory(namespace string, md5 string) string {
	return fmt.Sprintf("%s/%s", p.RemoteBasePath, p.imageDirectory(namespace, md5))
}

//  LocalImagePath returns local path for resized image
func (p *Paths) LocalImagePath(namespace string, md5 string, imageName string) string {
	return fmt.Sprintf("%s/%s", p.LocalBasePath, p.imagePath(namespace, md5, imageName))
}

func (p *Paths) RemoteImagePath(namespace string, md5 string, imageName string) string {
	return fmt.Sprintf("%s/%s", p.RemoteBasePath, p.imagePath(namespace, md5, imageName))
}

// RemoteOriginalPath returns local path for original image
func (p *Paths) RemoteOriginalPath(namespace string, md5 string) string {
	return fmt.Sprintf("%s/%s", p.RemoteBasePath, p.originalPath(namespace, md5))
}

func (p *Paths) RemoteOriginalURL(namespace string, md5 string) string {
	return fmt.Sprintf("%s/%s", p.RemoteBaseURL, p.RemoteOriginalPath(namespace, md5))
}

func (p *Paths) LocalInfoPath(namespace string, md5 string) string {
	return fmt.Sprintf("%s/%s", p.LocalBasePath, p.infoPath(namespace, md5))
}

func (p *Paths) RemoteInfoPath(namespace string, md5 string) string {
	return fmt.Sprintf("%s/%s", p.RemoteBasePath, p.infoPath(namespace, md5))
}

func (p *Paths) TempImagePath(url string) string {
	data := []byte(url)
	name := fmt.Sprintf("%x", md5.Sum(data))
	return fmt.Sprintf("%s/tmp/%s", p.LocalBasePath, name)
}

// originalPath
func (p *Paths) originalPath(namespace string, md5 string) string {
	return fmt.Sprintf("%s/original", p.imageDirectory(namespace, md5))
}

// imageDirectory returns relative directory starting at image root
func (p *Paths) imageDirectory(namespace string, md5 string) string {
	partitions := []string{md5[0:3], md5[3:6], md5[6:9], md5[9:32]}
	return fmt.Sprintf("%s/%s", namespace, strings.Join(partitions, "/"))
}

func (p *Paths) infoPath(namespace string, md5 string) string {
	return fmt.Sprintf("%s/info.json", p.imageDirectory(namespace, md5))
}

//  imagePath returns relative path to resized image
func (p *Paths) imagePath(namespace string, md5 string, imageName string) string {
	return fmt.Sprintf("%s/%s", p.imageDirectory(namespace, md5), imageName)
}
