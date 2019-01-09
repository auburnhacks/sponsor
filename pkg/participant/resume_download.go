package participant

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

var c = http.DefaultClient

// AllResumes is a function that lists all participants all downloads their
// resumes from google cloud and provides a byte sequence in tar format
func AllResumes() ([]byte, error) {
	participants, err := List()
	if err != nil {
		return nil, err
	}
	urls := []string{}
	for _, p := range participants {
		if len(p.Resume) == 0 {
			continue
		}
		urls = append(urls, p.Resume)
	}
	bt, err := Download(urls)
	if err != nil {
		return nil, err
	}
	var tbuf bytes.Buffer
	tw := tar.NewWriter(&tbuf)
	for i, bb := range bt {
		hdr := &tar.Header{
			Name: fmt.Sprintf("resume-%d.pdf", i),
			Mode: 0600,
			Size: int64(len(bb)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return nil, errors.Wrap(err,
				"pkg/participant: error while writing header.")
		}
		if _, err := tw.Write(bb); err != nil {
			return nil, errors.Wrap(err,
				"pkg/participant: error while wrint to tar archive")
		}
	}
	if err := tw.Close(); err != nil {
		return nil, errors.Wrap(err,
			"pkg/participant: error while closing tar archive")
	}
	return tbuf.Bytes(), nil
}

// Download ia function that will allow you to download all the resumes at once
func Download(urls []string) ([][]byte, error) {
	bufCh := make(chan []byte, len(urls))
	urlCh := make(chan string)

	for i := 0; i < len(urls); i++ {
		go fetch(urlCh, bufCh)
	}

	for _, url := range urls {
		urlCh <- url
	}
	close(urlCh)
	bt := [][]byte{} // bt is a bytes table that has all the resumes
	for i := 0; i < len(urls); i++ {
		bb := <-bufCh
		bt = append(bt, bb)
	}
	return bt, nil
}

func fetch(urlCh <-chan string, bufCh chan<- []byte) {
	url, ok := <-urlCh
	if !ok {
		return
	}
	bb, err := getBytes(url)
	if err != nil {
		return
	}
	bufCh <- bb
}

func getBytes(url string) ([]byte, error) {
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bb, nil
}
