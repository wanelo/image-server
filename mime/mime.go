package mime

import (
	"fmt"
	"strings"
)

func ExtToContentType(ext string) string {
	ext = strings.ToLower(ext)
	switch ext {
	case "jpg":
		return "image/jpeg"
	default:
		return fmt.Sprintf("image/%s", ext)
	}
}
