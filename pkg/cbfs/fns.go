package cbfs

import (
	"encoding/binary"
	"io"
	"log"
)

var Debug = log.Printf

func CBFSRead(r io.Reader, f interface{}) error {
	if err := binary.Read(r, CBFSEndian, f); err != nil {
		return err
	}
	return nil
}
