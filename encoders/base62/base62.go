/*
Package base62 encoding/decoding
Duplicate of https://bitbucket.org/tebeka/base62
Alphabet was modified to match ruby's base62 gem encoding/decoding
*/
package base62

import (
	"math"
	"strings"
)

const (
	Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	Version  = "0.2.0"
)

var base = uint64(len(Alphabet))

// Encode encodes num to base62 string.
func Encode(num uint64) string {
	if num == 0 {
		return "0"
	}

	arr := []uint8{}

	for num > 0 {
		rem := num % base
		num = num / base
		arr = append(arr, Alphabet[rem])
	}

	// Reverse the result array
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}

	return string(arr)
}

// Decode decodes base62 string to a number.
func Decode(b62 string) uint64 {
	size := len(b62)
	num := uint64(0)
	base := float64(len(Alphabet))

	for i, ch := range b62 {
		idx := i + 1
		loc := uint64(strings.IndexRune(Alphabet, ch))
		pow := uint64(math.Pow(base, float64(size-idx)))
		num += loc * pow
	}

	return num
}
