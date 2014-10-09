package noop

type Uploader struct{}

// Upload does nothing
func (u *Uploader) Upload(source string, destination string, contType string) error {
	return nil
}

// CreateDirectory does nothing
func (u *Uploader) CreateDirectory(path string) error {
	return nil
}

// ListDirectory does nothing and returns empty array
func (u *Uploader) ListDirectory(directory string) ([]string, error) {
	var names []string
	return names, nil
}
