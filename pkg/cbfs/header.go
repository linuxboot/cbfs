package cbfs

import "io"

func NewHeader(r io.Reader, f File) (CBFSReadWriter, error) {
	h := &CBFSFile{File: f}
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
