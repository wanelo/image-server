package main

import "testing"

func TestAllowedImageValid(t *testing.T) {
  fmts := []string{"jpg"}
  sc := &ServerConfiguration{MaximumWidth: 1000, WhitelistedExtensions: fmts}
  ic := &ImageConfiguration{width: 1000, format: "jpg"}

  allowed, _ := allowedImage(sc, ic)
  if allowed == false {
    t.Errorf("expected true")
  }
}

func TestAllowedImageTooWide(t *testing.T) {
  fmts := []string{"jpg"}
  sc := &ServerConfiguration{MaximumWidth: 1000, WhitelistedExtensions: fmts}
  ic := &ImageConfiguration{width: 1001, format: "jpg"}

  allowed, _ := allowedImage(sc, ic)
  if allowed == true {
    t.Errorf("expected false")
  }
}

func TestAllowedImageInvalidFormat(t *testing.T) {
  fmts := []string{"jpg"}
  sc := &ServerConfiguration{MaximumWidth: 1000, WhitelistedExtensions: fmts}
  ic := &ImageConfiguration{width: 100, format: "pdf"}

  allowed, _ := allowedImage(sc, ic)
  if allowed == true {
    t.Errorf("expected false")
  }
}
