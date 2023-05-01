package types

import (
	"bytes"
	"compress/zlib"
	"io"
)

func Decompress(data []byte) ([]byte, error) {
	// zip reader
	zr, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	return io.ReadAll(zr)
}
