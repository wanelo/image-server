package mime

import (
	"fmt"
	"strings"
)

// ExtToContentType returns the content type for a given file extension.
// The content type is retuned in the header when serving images
func ExtToContentType(ext string) string {
	ext = strings.ToLower(ext)
	switch ext {
	case "jpg":
		return "image/jpeg"
	case "pdf":
		return "application/pdf"
	default:
		return fmt.Sprintf("image/%s", ext)
	}
}
