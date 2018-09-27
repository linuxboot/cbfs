package cbfs

import (
	"log"
)

func init() {
	if err := RegisterFileReader(&SegReader{T: TypeBootBlock, N: "CBFSBootBlock", F: NewBootBlock}); err != nil {
		log.Fatal(err)
	}
}

func NewBootBlock(r CountingReader, f *File) (ReadWriter, error) {
	h := &BootBlockRecord{File: *f}
	Debug("Before BootBlock: total bytes read: %d", r.Count())
	Debug("Got header %v", *h)
	return h, nil
}

func (h *BootBlockRecord) Read([]byte) (int, error) {
	return -1, nil
}

func (h *BootBlockRecord) Write([]byte) (int, error) {
	return -1, nil
}

func (H *BootBlockRecord) String() string {
	return ""
}
