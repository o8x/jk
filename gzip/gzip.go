package gzip

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
)

func Encode(content []byte) (string, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	if _, err := zw.Write(content); err != nil {
		return "", err
	}

	if err := zw.Close(); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func Decode(text string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return nil, err
	}

	zr, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	if err := zr.Close(); err != nil {
		return nil, err
	}

	return io.ReadAll(zr)
}
