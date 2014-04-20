package cli

import (
	"reflect"
	"testing"

	"github.com/wanelo/image-server/core"
)

func TestFullSizeImage(t *testing.T) {
	sc := &core.ServerConfiguration{MaximumWidth: 1000, LocalBasePath: "public"}
	ic := &core.ImageConfiguration{ServerConfiguration: sc, Width: 0, Height: 0, Format: "jpg", Quality: 85, Model: "test", ImageType: "image", ID: "11032603", Filename: "full_size.jpg"}

	expected := []string{"-format", "jpg", "-flatten", "-background", "rgba\\(255,255,255,1\\)", "-quality", "85", "public/test/image/11032603/original", "public/test/image/11032603/full_size.jpg"}

	command := commandArgs(ic)
	if !reflect.DeepEqual(expected, command) {
		t.Errorf("expected\n%v got\n%v", expected, command)
	}
}

func TestImageWithWidth(t *testing.T) {
	sc := &core.ServerConfiguration{MaximumWidth: 1000, LocalBasePath: "public"}
	ic := &core.ImageConfiguration{ServerConfiguration: sc, Width: 600, Height: 0, Format: "jpg", Quality: 85, Model: "test", ImageType: "image", ID: "11032603", Filename: "w600.jpg"}

	expected := []string{"-format", "jpg", "-flatten", "-background", "rgba\\(255,255,255,1\\)", "-quality", "85", "-resize", "600", "public/test/image/11032603/original", "public/test/image/11032603/w600.jpg"}

	command := commandArgs(ic)
	if !reflect.DeepEqual(expected, command) {
		t.Errorf("expected\n%v got\n%v", expected, command)
	}
}

func TestImageWithWidthAndHeight(t *testing.T) {
	sc := &core.ServerConfiguration{MaximumWidth: 1000, LocalBasePath: "public"}
	ic := &core.ImageConfiguration{ServerConfiguration: sc, Width: 600, Height: 500, Format: "jpg", Quality: 85, Model: "test", ImageType: "image", ID: "11032603", Filename: "600x500.jpg"}

	expected := []string{"-format", "jpg", "-flatten", "-background", "rgba\\(255,255,255,1\\)", "-quality", "85", "-extent", "600x500", "public/test/image/11032603/original", "public/test/image/11032603/600x500.jpg"}

	command := commandArgs(ic)
	if !reflect.DeepEqual(expected, command) {
		t.Errorf("expected\n%v got\n%v", expected, command)
	}
}
