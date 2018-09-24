package cbfs

import (
	"io"
	"log"
)

func init() {
	if err := RegisterFileReader(&SegReader{T: 2, N: "CBFSHeader", F: NewHeader}); err != nil {
		log.Fatal(err)
	}
}

func NewHeader(r io.Reader) (CBFSReadWriter, error) {
	h := &CBFSHeader{}
	if err := CBFSRead(r, h); err != nil {
		return nil, err
	}
	return h, nil
}

func (h *CBFSHeader) Read([]byte) (int, error) {
	return -1, nil
}

func (h *CBFSHeader) Write([]byte) (int, error) {
	return -1, nil
}

func (H *CBFSHeader) String() string {
	return ""
}
