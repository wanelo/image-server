package server

// import (
// 	"testing"
//
// 	"github.com/wanelo/image-server/core"
// )
//
// func TestAllowedImageValid(t *testing.T) {
// 	fmts := []string{"jpg"}
// 	sc := &core.ServerConfiguration{MaximumWidth: 1000, WhitelistedExtensions: fmts}
// 	ic := &core.ImageConfiguration{ServerConfiguration: sc, Width: 1000, Format: "jpg"}
//
// 	allowed, _ := allowedImage(ic)
// 	if allowed == false {
// 		t.Errorf("expected true")
// 	}
// }
//
// func TestAllowedImageTooWide(t *testing.T) {
// 	fmts := []string{"jpg"}
// 	sc := &core.ServerConfiguration{MaximumWidth: 1000, WhitelistedExtensions: fmts}
// 	ic := &core.ImageConfiguration{ServerConfiguration: sc, Width: 1001, Format: "jpg"}
//
// 	allowed, _ := allowedImage(ic)
// 	if allowed == true {
// 		t.Errorf("expected false")
// 	}
// }
//
// func TestAllowedImageInvalidFormat(t *testing.T) {
// 	fmts := []string{"jpg"}
// 	sc := &core.ServerConfiguration{MaximumWidth: 1000, WhitelistedExtensions: fmts}
// 	ic := &core.ImageConfiguration{ServerConfiguration: sc, Width: 100, Format: "pdf"}
//
// 	allowed, _ := allowedImage(ic)
// 	if allowed == true {
// 		t.Errorf("expected false")
// 	}
// }
