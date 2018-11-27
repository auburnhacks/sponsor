package participant

import (
	"archive/tar"
	"bytes"
	"fmt"

	"github.com/pkg/errors"
)

// BulkDownload is a function that will allow authorized sponsors to download
// all resumes in a single shot for offline viewing.
func BulkDownload() ([]byte, error) {
	// Attempt to fetch all the resumes concurrently
	// Run it through the tar function
	// Run it through gzip function to minimize the latency while downloading
	pSlice, err := List()
	if err != nil {
		return nil, err
	}
	rBufCh := make(chan bytes.Buffer)
	errCh := make(chan error)
	for _, p := range pSlice {
		go fetchBytes(p.Resume, rBufCh, errCh)
	}
	var tBuf bytes.Buffer
	resIdx := 1
	tw := tar.NewWriter(&tBuf)
	for rBuf := range rBufCh {
		hdr := &tar.Header{
			Name: fmt.Sprintf("resume-%d,pdf", resIdx),
			Mode: 0600,
			Size: int64(rBuf.Len()),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return nil, err
		}
		if _, err := tw.Write(rBuf.Bytes()); err != nil {
			return nil, errors.Wrap(err, "pkg/participant: error while write to the tar archive")
		}
	}
	if err := tw.Close(); err != nil {
		return nil, errors.Wrap(err, "pkg/participant: error while closing tar archive")
	}
	return tBuf.Bytes(), nil
}

func fetchBytes(url string, outChan chan<- bytes.Buffer, errCh chan<- error) {
}
