// SPDX-License-Identifier: Apache-2.0

package util

import (
	"bytes"
	"compress/zlib"
	"io"
)

// Compress is a helper function to compress values. First, an
// empty buffer is created for storing compressed data. Then,
// a zlib writer, using the DEFLATE algorithm, is created with
// the provided compression level to output to this buffer.
// Finally, the provided value is compressed and written to the
// buffer and the writer is closed which flushes all bytes from
// the writer to the buffer.
func Compress(level int, value []byte) ([]byte, error) {
	// create new buffer for storing compressed data
	b := new(bytes.Buffer)

	// create new zlib writer for outputting data to the buffer in a compressed format
	w, err := zlib.NewWriterLevel(b, level)
	if err != nil {
		return value, err
	}

	// write data to the buffer in compressed format
	_, err = w.Write(value)
	if err != nil {
		return value, err
	}

	// close the writer
	//
	// compressed bytes are not flushed until the writer is closed or explicitly flushed
	err = w.Close()
	if err != nil {
		return value, err
	}

	// return compressed bytes from the buffer
	return b.Bytes(), nil
}

// Decompress is a helper function to decompress values. First, a
// buffer is created from the provided compressed data. Then, a
// zlib reader, using the DEFLATE algorithm, is created from the
// buffer as an input for reading data from the buffer. Finally,
// the data is decompressed and read from the buffer.
func Decompress(value []byte) ([]byte, error) {
	// create new buffer from the compressed data
	b := bytes.NewBuffer(value)

	// create new zlib reader for reading the compressed data from the buffer
	r, err := zlib.NewReader(b)
	if err != nil {
		return value, err
	}

	// close the reader after the data has been decompressed
	defer r.Close()

	// capture decompressed data from the compressed data in the buffer
	data, err := io.ReadAll(r)
	if err != nil {
		return value, err
	}

	return data, nil
}
