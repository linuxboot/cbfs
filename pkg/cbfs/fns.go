package cbfs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

var Debug = log.Printf

// Read reads things in in BE format, which they are supposed to be in.
func Read(r io.Reader, f interface{}) error {
	if err := binary.Read(r, Endian, f); err != nil {
		return err
	}
	return nil
}

// Read LE reads things in LE format, which the spec says it is not in.
func ReadLE(r io.Reader, f interface{}) error {
	if err := binary.Read(r, binary.LittleEndian, f); err != nil {
		return err
	}
	return nil
}

func ReadName(r CountingReader, f *File, size uint32) (string, error) {
	b := make([]byte, size)
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
