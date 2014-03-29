package main

import (
	"fmt"
	"regexp"
	"strconv"
)

func NameToConfiguration(filename string) (*ImageConfiguration, error) {
	var w, h string

	reR := regexp.MustCompile(`^([0-9]+)x([0-9]+)$`)
	reS := regexp.MustCompile(`^x([0-9]+)$`)
	reW := regexp.MustCompile(`^w([0-9]+)$`)
	reF := regexp.MustCompile(`^full_size$`)

	if reR.MatchString(filename) {
		m := reR.FindStringSubmatch(filename)
		w, h = m[1], m[2]
	} else if reS.MatchString(filename) {
		m := reS.FindStringSubmatch(filename)
		w, h = m[1], m[1]
	} else if reW.MatchString(filename) {
		m := reW.FindStringSubmatch(filename)
		w, h = m[1], "0"
	} else if reF.MatchString(filename) {
		w, h = "0", "0"
	} else {
		// return error
		return &ImageConfiguration{}, fmt.Errorf("Unsupported")
	}

	width, _ := strconv.Atoi(w)
	height, _ := strconv.Atoi(h)
	return &ImageConfiguration{width: width, height: height}, nil
}
