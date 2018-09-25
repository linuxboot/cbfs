package cbfs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

var Debug = log.Printf

func Read(r io.Reader, f interface{}) error {
	if err := binary.Read(r, Endian, f); err != nil {
		return err
	}
	return nil
}

func ReadName(r CountingReader, f *File) (string, error) {
	Debug("Readname: f.Offset %d, count %d", f.Offset, f.Offset - 24)
	b := make([]byte, f.Offset - 24)
	n, err := r.Read(b)
	if err != nil {
		Debug("ReadName failed:%v", err)
		return "", err
	}
	Debug("Readname gets %#02x", b)
	if n != len(b) {
		err = fmt.Errorf("ReadName: got %d, want %d for name", n, len(b))
		Debug("Readname short: %v", err)
		return "", err
	}
	// discard trailing NULLs
	z := bytes.Split(b, []byte{0})
	return string(z[0]), nil
}

func Align(r CountingReader) error {
	var junk [16]byte
	align := (int(r.Count()) + 15) & ^0xf
	amt := align - int(r.Count())
	return Read(r, junk[:amt])
}

