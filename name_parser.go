package main

import (
	"fmt"
	"regexp"
	"strconv"
)

func NameToConfiguration(filename string) (*ImageConfiguration, error) {
	var w, h, f string

	reR := regexp.MustCompile(`^([0-9]+)x([0-9]+)\.(\w{3,4})$`)
	reS := regexp.MustCompile(`^x([0-9]+)\.(\w{3,4})$`)
	reW := regexp.MustCompile(`^w([0-9]+)\.(\w{3,4})$`)
	reF := regexp.MustCompile(`^full_size\.(\w{3,4})$`)

	if reR.MatchString(filename) {
		m := reR.FindStringSubmatch(filename)
		w, h, f = m[1], m[2], m[3]
	} else if reS.MatchString(filename) {
		m := reS.FindStringSubmatch(filename)
		w, h, f = m[1], m[1], m[2]
	} else if reW.MatchString(filename) {
		m := reW.FindStringSubmatch(filename)
		w, h, f = m[1], "0", m[2]
	} else if reF.MatchString(filename) {
		m := reF.FindStringSubmatch(filename)
		w, h, f = "0", "0", m[1]
	} else {
		// return error
		return &ImageConfiguration{}, fmt.Errorf("Unsupported")
	}

	width, _ := strconv.Atoi(w)
	height, _ := strconv.Atoi(h)
	return &ImageConfiguration{width: width, height: height, format: f}, nil
}
