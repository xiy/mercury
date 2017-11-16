package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
)

// Compress compreses a string with the Gzip algorithm
func Compress(s string) string {
	var b bytes.Buffer

	gz, err := gzip.NewWriterLevel(&b, gzip.BestSpeed)

	_, err = gz.Write([]byte(s))
	if err != nil {
		panic(err)
	}

	if err = gz.Close(); err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(b.Bytes())
}

// Decompress decompresses a gzip compressed string
func Decompress(s string) string {
	// Write gzipped data to the client
	decodedString, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		logger.Fatalf("error: %s", err)
	}

	buf := bytes.NewBuffer(decodedString)

	gr, err := gzip.NewReader(buf)
	if err != nil {
		logger.Fatalf("error: %s", err)
	}
	defer gr.Close()

	data, err := ioutil.ReadAll(gr)
	if err != nil {
		logger.Fatalf("error: %s", err)
	}

	return string(data)
}
