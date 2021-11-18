package io

import "io"

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

func (a nopCloser) Seek(offset int64, whence int) (int64, error)  {
	seeker := a.Reader.(io.Seeker)
	return seeker.Seek(offset, whence)
}

func (nopCloser) Close() error { return nil }
