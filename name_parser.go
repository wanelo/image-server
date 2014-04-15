package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/wanelo/image-server/core"
)

func NameToConfiguration(sc *core.ServerConfiguration, filename string) (*core.ImageConfiguration, error) {
	var w, h, q, f string
	var quality uint

	reR := regexp.MustCompile(`^([0-9]+)x([0-9]+)(?:-q([0-9]+))?\.(\w{3,4})$`)
	reS := regexp.MustCompile(`^x([0-9]+)(?:-q([0-9]+))?\.(\w{3,4})$`)
	reW := regexp.MustCompile(`^w([0-9]+)(?:-q([0-9]+))?\.(\w{3,4})$`)
	reF := regexp.MustCompile(`^full_size(?:-q([0-9]+))?\.(\w{3,4})$`)

	if reR.MatchString(filename) {
		m := reR.FindStringSubmatch(filename)
		w, h, q, f = m[1], m[2], m[3], m[4]
	} else if reS.MatchString(filename) {
		m := reS.FindStringSubmatch(filename)
		w, h, q, f = m[1], m[1], m[2], m[3]
	} else if reW.MatchString(filename) {
		m := reW.FindStringSubmatch(filename)
		w, h, q, f = m[1], "0", m[2], m[3]
	} else if reF.MatchString(filename) {
		m := reF.FindStringSubmatch(filename)
		w, h, q, f = "0", "0", m[1], m[2]
	} else {
		// return error
		return &core.ImageConfiguration{}, fmt.Errorf("unsupported")
	}

	width, _ := strconv.Atoi(w)
	height, _ := strconv.Atoi(h)
	quality64, _ := strconv.ParseUint(q, 10, 0)

	if quality64 > 0 {
		quality = uint(quality64)
	} else {
		quality = sc.DefaultQuality
	}

	return &core.ImageConfiguration{Width: width, Height: height, Quality: quality, Format: f, Filename: filename}, nil
}
