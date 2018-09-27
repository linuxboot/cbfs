package cbfs

import (
	"log"
)

func init() {
	if err := RegisterFileReader(&SegReader{T: TypeRaw, N: "CBFSRaw", F: NewRaw}); err != nil {
		log.Fatal(err)
	}
}

func NewRaw(r CountingReader, f *File) (ReadWriter, error) {
	h := &RawRecord{File: *f}
	Debug("Before Raw: total bytes read: %d", r.Count())
	Debug("Got header %v", *h)
	return h, nil
}

func (h *RawRecord) Read([]byte) (int, error) {
	return -1, nil
}

func (h *RawRecord) Write([]byte) (int, error) {
	return -1, nil
}

func (H *RawRecord) String() string {
	return ""
}