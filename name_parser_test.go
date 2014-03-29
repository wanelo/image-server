package main

import "testing"

func ensureDimensions(t *testing.T, ic *ImageConfiguration, w int, h int, f string) {
	if ic.width != w {
		t.Errorf("expected %v to be %v", ic.width, w)

	}
	if ic.height != h {
		t.Errorf("expected %v to be %v", ic.width, h)
	}
	if ic.format != f {
		t.Errorf("expected %v to be %v", ic.format, f)
	}
}

func TestRectangle(t *testing.T) {
	ic, _ := NameToConfiguration("300x400.jpg")
	ensureDimensions(t, ic, 300, 400, "jpg")
}

func TestSquare(t *testing.T) {
	ic, _ := NameToConfiguration("x300.jpg")
	ensureDimensions(t, ic, 300, 300, "jpg")
}

func TestWidth(t *testing.T) {
	ic, _ := NameToConfiguration("w300.jpg")
	ensureDimensions(t, ic, 300, 0, "jpg")
}

func TestFullSize(t *testing.T) {
	ic, _ := NameToConfiguration("full_size.jpg")
	ensureDimensions(t, ic, 0, 0, "jpg")
}

func TestUnsupported(t *testing.T) {
	_, err := NameToConfiguration("random.jpg")
	if err == nil {
		t.Errorf("expected to receive an error")
	}
}
