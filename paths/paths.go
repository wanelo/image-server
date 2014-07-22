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
}

// OriginalPath
func (p *Paths) OriginalPath(namespace string, md5 string) string {
	return fmt.Sprintf("%s/original", p.ImageDirectory(namespace, md5))
}

// ImageDirectory returns relative directory starting at image root
func (p *Paths) ImageDirectory(namespace string, md5 string) string {
	partitions := []string{md5[0:3], md5[3:6], md5[6:9], md5[9:32]}
	return fmt.Sprintf("%s/%s", namespace, strings.Join(partitions, "/"))
}

//  ImagePath returns relative path to resized image
func (p *Paths) ImagePath(namespace string, md5 string, imageName string) string {
	return fmt.Sprintf("%s/%s", p.ImageDirectory(namespace, md5), imageName)
}

// LocalOriginalPath returns local path for original image
func (p *Paths) LocalOriginalPath(namespace string, md5 string) string {
	return fmt.Sprintf("%s/%s", p.LocalBasePath, p.OriginalPath(namespace, md5))
}

// LocalImageDirectory returns location for locally cached images
func (p *Paths) LocalImageDirectory(namespace string, md5 string) string {
	return fmt.Sprintf("%s/%s", p.LocalBasePath, p.ImageDirectory(namespace, md5))
}

//  LocalImagePath returns local path for resized image
func (p *Paths) LocalImagePath(namespace string, md5 string, imageName string) string {
	return fmt.Sprintf("%s/%s", p.LocalBasePath, p.ImagePath(namespace, md5, imageName))
}

// RemoteOriginalPath returns local path for original image
func (p *Paths) RemoteOriginalPath(namespace string, md5 string) string {
	return fmt.Sprintf("%s/%s", p.RemoteBasePath, p.OriginalPath(namespace, md5))
}

func (p *Paths) InfoPath(namespace string, md5 string) string {
	return fmt.Sprintf("%s/info.json", p.ImageDirectory(namespace, md5))
}

func (p *Paths) LocalInfoPath(namespace string, md5 string) string {
	return fmt.Sprintf("%s/%s", p.LocalBasePath, p.InfoPath(namespace, md5))
}

func (p *Paths) RemoteInfoPath(namespace string, md5 string) string {
	return fmt.Sprintf("%s/%s", p.RemoteBasePath, p.InfoPath(namespace, md5))
}

func (p *Paths) TempDirectory() string {
	return fmt.Sprintf("%s/tmp", p.LocalBasePath)
}

func (p *Paths) TempImagePath(url string) string {
	data := []byte(url)
	name := fmt.Sprintf("%x", md5.Sum(data))
	return fmt.Sprintf("%s/%s", p.TempDirectory(), name)
}
