package cbfs

import (
	"log"
)

func init() {
	if err := RegisterFileReader(&SegReader{T: TypeDeleted, N: "CBFSEmpty", F: NewEmpty}); err != nil {
		log.Fatal(err)
	}
	if err := RegisterFileReader(&SegReader{T: TypeDeleted2, N: "CBFSEmpty", F: NewEmpty}); err != nil {
		log.Fatal(err)
	}
}

func NewEmpty(r CountingReader, f *File) (ReadWriter, error) {
	h := &EmptyRecord{File: *f}
	Debug("Before Empty: total bytes read: %d", r.Count())
	Debug("Got header %v", *h)
	h.Data = make([]byte, h.Size)
	n, err := r.Read(h.Data)
	if err != nil {
		return nil, err
	}
	Debug("Bootblock read %d bytes, now at %#x", n, r.Count())

	return h, nil
}

func (h *EmptyRecord) Read([]byte) (int, error) {
	return -1, nil
}

func (h *EmptyRecord) Write([]byte) (int, error) {
	return -1, nil
}

func (h *EmptyRecord) String() string {
	return recString("(empty)", h.RomOffset, h.Type.String(), h.Size, "none")
}

func (h *EmptyRecord) Name() string {
	return "(empty)"
}
