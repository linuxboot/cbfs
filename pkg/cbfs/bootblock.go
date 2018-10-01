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
	h.Data = make([]byte, h.Size)
	n, err := r.Read(h.Data)
	if err != nil {
		return nil, err
	}
	Debug("Bootblock read %d bytes", n)
	return h, nil
}

func (h *BootBlockRecord) Read([]byte) (int, error) {
	return -1, nil
}

func (h *BootBlockRecord) Write([]byte) (int, error) {
	return -1, nil
}

func (h *BootBlockRecord) String() string {
	return recString(h.Name(), h.RomOffset, h.Type.String(), h.Size, "none")
}

func (h *BootBlockRecord) Name() string {
	return "BootBlock"
}
