package cache

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"strings"
)

// URLHash returns an MD5 hash of a normalized URL
func URLHash(u *url.URL) string {
	return getMD5Hash(normalizeURL(*u))
}

// removew URL scheme and makes it case insensitive
func normalizeURL(u url.URL) string {
	u.Scheme = ""
	return strings.ToLower(u.String())
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
