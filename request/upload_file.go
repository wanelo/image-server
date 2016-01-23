package request

import (
	"io"
	"os"

	"github.com/image-server/image-server/mime"
	"github.com/image-server/image-server/uploader"
)

func (r *Request) UploadFile(filename string) error {
	localDirectory := r.Paths.LocalImageDirectory(r.Namespace, r.Hash)
	os.MkdirAll(localDirectory, 0700)

	localPath := r.Paths.LocalImagePath(r.Namespace, r.Hash, filename)

	out, err := os.Create(localPath)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, r.SourceData)
	if err != nil {
		return err
	}

	uploader := uploader.DefaultUploader(r.ServerConfiguration)
	err = uploader.CreateDirectory(r.Paths.RemoteImageDirectory(r.Namespace, r.Hash))
	if err != nil {
		return err
	}

	remotePath := r.Paths.RemoteImagePath(r.Namespace, r.Hash, filename)

	contentType := mime.ExtToContentType(filename)
	// upload original image
	err = uploader.Upload(localPath, remotePath, contentType)
	if err != nil {
		return err
	}
	return nil
}
