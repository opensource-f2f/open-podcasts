package io

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// SeekableReader represents a reader that be able to seek
type SeekableReader interface {
	io.Reader
	io.Seeker
	io.ReadCloser
}

// SeekerWithoutCloser creates a NopReader has seeking feature
func SeekerWithoutCloser(r io.Reader) SeekableReader {
	return nopCloser{r}
}

type nopCloser struct {
	io.Reader
}

func (a nopCloser) Seek(offset int64, whence int) (int64, error) {
	seeker := a.Reader.(io.Seeker)
	return seeker.Seek(offset, whence)
}

func (nopCloser) Close() error { return nil }

type RangeReader struct {
	offset     int
	length     int
	bufferSize int

	resource string
}

func NewRangeReader(offset, lenght int, resource string) *RangeReader {
	return &RangeReader{
		offset:     offset,
		length:     lenght,
		bufferSize: 0,
		resource:   resource,
	}
}

func (r *RangeReader) Read(p []byte) (n int, err error) {
	if r.bufferSize == 0 {
		r.bufferSize = len(p)
	}

	client := http.DefaultClient

	var request *http.Request
	if request, err = http.NewRequest(http.MethodGet, r.resource, nil); err != nil {
		return
	}

	request.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", r.offset, r.offset+r.bufferSize))

	var resp *http.Response
	if resp, err = client.Do(request); err != nil {
		return
	}

	if resp.StatusCode == http.StatusPartialContent {
		if n, err = resp.Body.Read(p); err == nil {
			r.offset += r.bufferSize
		}
	} else {
		return 0, errors.New(fmt.Sprintf("unspport status code: %d", resp.StatusCode))
	}
	return
}
