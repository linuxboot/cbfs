package cbfs

import (
	"log"
)

func init() {
	if err := RegisterFileReader(&SegReader{T: TypeCMOSLayout, N: "CBFSCMOSLayout", F: NewCMOSLayout}); err != nil {
		log.Fatal(err)
	}
}

func NewCMOSLayout(r CountingReader, f *File) (ReadWriter, error) {
	h := &CMOSLayoutRecord{File: *f}
	Debug("Before CMOSLayout: total bytes read: %d", r.Count())
	Debug("Got header %v", *h)
	return h, nil
}

func (h *CMOSLayoutRecord) Read([]byte) (int, error) {
	return -1, nil
}

func (h *CMOSLayoutRecord) Write([]byte) (int, error) {
	return -1, nil
}

func (H *CMOSLayoutRecord) String() string {
	return ""
}
