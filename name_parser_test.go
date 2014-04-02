package main

import "testing"

func ensureImageConfiguration(t *testing.T, ic *ImageConfiguration, w int, h int, q uint, f string) {
	if ic.width != w {
		t.Errorf("expected %v to be %v", ic.width, w)

	}
	if ic.height != h {
		t.Errorf("expected %v to be %v", ic.width, h)
	}
	if ic.quality != q {
		t.Errorf("expected %v to be %v", ic.quality, q)
	}
	if ic.format != f {
		t.Errorf("expected %v to be %v", ic.format, f)
	}
}

var sc *ServerConfiguration

func init() {
	sc, _ = loadServerConfiguration("test")
}

// Use the default quality

func TestRectangle(t *testing.T) {
	ic, _ := NameToConfiguration(sc, "300x400.jpg")
	ensureImageConfiguration(t, ic, 300, 400, 75, "jpg")
}

func TestSquare(t *testing.T) {
	ic, _ := NameToConfiguration(sc, "x300.jpg")
	ensureImageConfiguration(t, ic, 300, 300, 75, "jpg")
}

func TestWidth(t *testing.T) {
	ic, _ := NameToConfiguration(sc, "w300.jpg")
	ensureImageConfiguration(t, ic, 300, 0, 75, "jpg")
}

func TestFullSize(t *testing.T) {
	ic, _ := NameToConfiguration(sc, "full_size.jpg")
	ensureImageConfiguration(t, ic, 0, 0, 75, "jpg")
}

func TestUnsupported(t *testing.T) {
	_, err := NameToConfiguration(sc, "random.jpg")
	if err == nil {
		t.Errorf("expected to receive an error")
	}
}

// Quality is Provided

func TestRectangleWithQuality(t *testing.T) {
	ic, _ := NameToConfiguration(sc, "300x400-q10.jpg")
	ensureImageConfiguration(t, ic, 300, 400, 10, "jpg")
}

func TestSquareWithQuality(t *testing.T) {
	ic, _ := NameToConfiguration(sc, "x300-q10.jpg")
	ensureImageConfiguration(t, ic, 300, 300, 10, "jpg")
}

func TestWidthWithQuality(t *testing.T) {
	ic, _ := NameToConfiguration(sc, "w300-q10.jpg")
	ensureImageConfiguration(t, ic, 300, 0, 10, "jpg")
}

func TestFullSizeWithQuality(t *testing.T) {
	ic, _ := NameToConfiguration(sc, "full_size-q10.jpg")
	ensureImageConfiguration(t, ic, 0, 0, 10, "jpg")
}
