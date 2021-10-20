package validator

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func gzipDecompress(in []byte) []byte {
	reader := bytes.NewReader(in)
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return in
	}
	output, err := ioutil.ReadAll(gzReader)
	if err != nil {
		return in
	}
	return output
}
