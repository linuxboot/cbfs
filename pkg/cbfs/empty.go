package cbfs

import (
	"io"
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

func (r *EmptyRecord) Read([]byte) (int, error) {
	return -1, nil
}

func (r *EmptyRecord) Write([]byte) (int, error) {
	return -1, nil
}

func (r *EmptyRecord) String() string {
	return recString("(empty)", r.RomOffset, r.Type.String(), r.Size, "none")
}

func (r *EmptyRecord) Name() string {
	return "(empty)"
}

func (r *EmptyRecord) Update(w io.Writer) error {
	if err := Write(w, r.FileHeader); err != nil {
		return err
	}
	return Write(w, r.Data)
}

func (r *EmptyRecord) Header() *File {
	return &r.File
}
