package buff

import (
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
)

type Compression int8

const (
	None Compression = iota
	Gzip
	Zlib
)

func ArchiveReader(kind Compression, r io.Reader) (a io.Reader) {
	var err error

	switch kind {
	case None:
		a = r
	case Gzip:
		a, err = gzip.NewReader(r)
	case Zlib:
		a, err = zlib.NewReader(r)
	default:
		err = fmt.Errorf("Unknown Compression Type %d", kind)
	}

	if err != nil {
		panic(err)
	}

	return a
}

func ArchiveWriter(kind Compression, w io.Writer) (a io.Writer) {
	var err error

	switch kind {
	case None:
		a = w
	case Gzip:
		a = gzip.NewWriter(w)
	case Zlib:
		a = zlib.NewWriter(w)
	default:
		err = fmt.Errorf("Unknown Compression Type %d", kind)
	}

	if err != nil {
		panic(err)
	}

	return a
}
