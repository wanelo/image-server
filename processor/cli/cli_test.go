package cli

import (
	"reflect"
	"testing"

	"github.com/wanelo/image-server/core"
)

func TestFullSizeImage(t *testing.T) {
	sc := &core.ServerConfiguration{MaximumWidth: 1000, LocalBasePath: "public"}
	ic := &core.ImageConfiguration{ServerConfiguration: sc, Width: 0, Height: 0, Format: "jpg", Quality: 85, Namespace: "test", ID: "ofrA", Filename: "full_size.jpg"}

	expected := []string{"-format", "jpg", "-flatten", "-background", "rgba\\(255,255,255,1\\)", "-quality", "85", "public/test/00/of/rA/original", "public/test/00/of/rA/full_size.jpg"}

	command := commandArgs(ic)
	if !reflect.DeepEqual(expected, command) {
		t.Errorf("expected\n%v got\n%v", expected, command)
	}
}

func TestImageWithWidth(t *testing.T) {
	sc := &core.ServerConfiguration{MaximumWidth: 1000, LocalBasePath: "public"}
	ic := &core.ImageConfiguration{ServerConfiguration: sc, Width: 600, Height: 0, Format: "jpg", Quality: 85, Namespace: "test", ID: "ofrA", Filename: "w600.jpg"}

	expected := []string{"-format", "jpg", "-flatten", "-background", "rgba\\(255,255,255,1\\)", "-quality", "85", "-resize", "600", "public/test/00/of/rA/original", "public/test/00/of/rA/w600.jpg"}

	command := commandArgs(ic)
	if !reflect.DeepEqual(expected, command) {
		t.Errorf("expected\n%v got\n%v", expected, command)
	}
}

func TestImageWithWidthAndHeight(t *testing.T) {
	sc := &core.ServerConfiguration{MaximumWidth: 1000, LocalBasePath: "public"}
	ic := &core.ImageConfiguration{ServerConfiguration: sc, Width: 600, Height: 500, Format: "jpg", Quality: 85, Namespace: "test", ID: "ofrA", Filename: "600x500.jpg"}

	expected := []string{"-format", "jpg", "-flatten", "-background", "rgba\\(255,255,255,1\\)", "-quality", "85", "-extent", "600x500", "-gravity", "center", "public/test/00/of/rA/original", "public/test/00/of/rA/600x500.jpg"}

	command := commandArgs(ic)
	if !reflect.DeepEqual(expected, command) {
		t.Errorf("expected\n%v got\n%v", expected, command)
	}
}
