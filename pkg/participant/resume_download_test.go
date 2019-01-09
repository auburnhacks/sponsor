package participant

import (
	"testing"
)

func TestGetBytes(t *testing.T) {
	url := "https://auburnhacks.com/metadata"
	_, err := getBytes(url)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestDownload(t *testing.T) {
	urls := []string{
		"https://storage.googleapis.com/resumes_19/54708071-4f77-4f62-9bed-6a903e4248c9.pdf",
		"https://storage.googleapis.com/resumes_19/54708071-4f77-4f62-9bed-6a903e4248c9.pdf",
		"https://storage.googleapis.com/resumes_19/54708071-4f77-4f62-9bed-6a903e4248c9.pdf",
		"https://storage.googleapis.com/resumes_19/54708071-4f77-4f62-9bed-6a903e4248c9.pdf",
	}
	_, err := Download(urls)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}
