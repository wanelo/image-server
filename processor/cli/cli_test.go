package cli

import (
	"reflect"
	"testing"

	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/info"
)

func TestFullSizeImage(t *testing.T) {
	ic := &core.ImageConfiguration{Width: 0, Height: 0, Format: "jpg", Quality: 85, Namespace: "test", ID: "ofrA", Filename: "full_size.jpg"}

	expected := []string{"-format", "jpg", "-flatten", "-background", "rgba(255,255,255,1)", "-quality", "85", "public/test/00/of/rA/original", "public/test/00/of/rA/full_size.jpg"}
	p := Processor{
		Source:             "public/test/00/of/rA/original",
		Destination:        "public/test/00/of/rA/full_size.jpg",
		ImageConfiguration: ic,
	}
	command := p.commandArgs()
	if !reflect.DeepEqual(expected, command) {
		t.Errorf("expected\n%v got\n%v", expected, command)
	}
}

func TestImageWithWidth(t *testing.T) {
	ic := &core.ImageConfiguration{Width: 600, Height: 0, Format: "jpg", Quality: 85, Namespace: "test", ID: "ofrA", Filename: "w600.jpg"}

	expected := []string{"-format", "jpg", "-flatten", "-resize", "600", "-background", "rgba(255,255,255,1)", "-quality", "85", "public/test/00/of/rA/original", "public/test/00/of/rA/w600.jpg"}

	p := Processor{
		Source:             "public/test/00/of/rA/original",
		Destination:        "public/test/00/of/rA/w600.jpg",
		ImageConfiguration: ic,
	}
	command := p.commandArgs()
	if !reflect.DeepEqual(expected, command) {
		t.Errorf("expected\n%v got\n%v", expected, command)
	}
}

func TestImageWithWidthAndHeight(t *testing.T) {
	ic := &core.ImageConfiguration{Width: 600, Height: 500, Format: "jpg", Quality: 85, Namespace: "test", ID: "ofrA", Filename: "600x500.jpg"}
	id := &info.ImageDetails{Width: 600, Height: 500}

	expected := []string{"-format", "jpg", "-flatten", "-extent", "600x500", "-gravity", "center", "-background", "rgba(255,255,255,1)", "-quality", "85", "public/test/00/of/rA/original", "public/test/00/of/rA/600x500.jpg"}

	p := Processor{
		Source:             "public/test/00/of/rA/original",
		Destination:        "public/test/00/of/rA/600x500.jpg",
		ImageConfiguration: ic,
		ImageDetails:       id,
	}
	command := p.commandArgs()
	if !reflect.DeepEqual(expected, command) {
		t.Errorf("expected\n%v got\n%v", expected, command)
	}
}
