package request

import (
	"github.com/golang/glog"
	"github.com/image-server/image-server/core"
	"github.com/image-server/image-server/info"
	"github.com/image-server/image-server/processor"
)

func (r *Request) Process(ic *core.ImageConfiguration) error {
	// The original file will be downloaded only once, even when every dimension requests it
	err := r.DownloadOriginal()
	if err != nil {
		return err
	}

	localResizedPath := r.Paths.LocalImagePath(r.Namespace, r.Hash, ic.Filename)
	localOriginalPath := r.Paths.LocalOriginalPath(r.Namespace, r.Hash)

	// process image
	pchan := &processor.ProcessorChannels{
		ImageProcessed: make(chan *core.ImageConfiguration),
		Skipped:        make(chan string),
	}

	info := &info.Info{
		Path: localOriginalPath,
	}
	id, err := info.ImageDetails()
	if err != nil {
		return err
	}

	p := processor.Processor{
		Source:             localOriginalPath,
		Destination:        localResizedPath,
		ImageConfiguration: ic,
		ImageDetails:       id,
		Channels:           pchan,
	}

	err = p.CreateImage()
	if err != nil {
		return err
	}

	select {
	case <-pchan.ImageProcessed:
		glog.Infof("Processed (resize handler) %s", localResizedPath)
		go r.uploadResizedImage(localResizedPath, ic)

	case <-pchan.Skipped:
		glog.Infof("Skipped processing (resize handler) %s", localResizedPath)
	}

	return nil
}

func (r *Request) uploadResizedImage(localResizedPath string, ic *core.ImageConfiguration) (err error) {
	remoteResizedPath := r.Paths.RemoteImagePath(ic.Namespace, ic.ID, ic.Filename)
	err = r.Uploader.Upload(localResizedPath, remoteResizedPath, ic.ToContentType())

	if err != nil {
		return err
	}
	return nil
}
