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

func NewHeader(r io.Reader, f *File) (CBFSReadWriter, error) {
	h := &CBFSMasterRecord{File: *f}
	if err := CBFSRead(r, &h.CBFSHeader); err != nil {
		Debug("Header read: %v", err)
		return nil, err
	}
	Debug("Got header %v", *h)
	return h, nil
}

func (h *CBFSMasterRecord) Read([]byte) (int, error) {
	return -1, nil
}

func (h *CBFSMasterRecord) Write([]byte) (int, error) {
	return -1, nil
}

func (H *CBFSMasterRecord) String() string {
	return ""
}
