package main

import "testing"

func ensureDimensions(t *testing.T, ic *ImageConfiguration, w int, h int) {
  if ic.width != w {
    t.Errorf("expected %v to be %v", ic.width, w)

  }
  if ic.height != h {
    t.Errorf("expected %v to be %v", ic.width, h)
  }
}

func TestRectangle(t *testing.T) {
  ic, _ := NameToConfiguration("300x400")
  ensureDimensions(t, ic, 300, 400)
}

func TestSquare(t *testing.T) {
  ic, _ := NameToConfiguration("x300")
  ensureDimensions(t, ic, 300, 300)
}

func TestWidth(t *testing.T) {
  ic, _ := NameToConfiguration("w300")
  ensureDimensions(t, ic, 300, 0)
}

func TestFullSize(t *testing.T) {
  ic, _ := NameToConfiguration("full_size")
  ensureDimensions(t, ic, 0, 0)
}

func TestUnsupported(t *testing.T) {
  _, err := NameToConfiguration("random")
  if err == nil {
    t.Errorf("expected to receive an error")
  }
}
