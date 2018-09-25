package cbfs

import (
	"log"
)

func init() {
	if err := RegisterFileReader(&SegReader{T: 2, N: "CBFSHeader", F: NewHeader}); err != nil {
		log.Fatal(err)
	}
}

func NewHeader(r CountingReader, f *File) (ReadWriter, error) {
	h := &MasterRecord{File: *f}
	Debug("Before Header: total bytes read: %d", r.Count())
	if err := Read(r, &h.MasterHeader); err != nil {
		Debug("Header read: %v", err)
		return nil, err
	}
	Debug("Got header %v", *h)
	return h, nil
}

func (h *MasterRecord) Read([]byte) (int, error) {
	return -1, nil
}

func (h *MasterRecord) Write([]byte) (int, error) {
	return -1, nil
}

func (H *MasterRecord) String() string {
	return ""
}
