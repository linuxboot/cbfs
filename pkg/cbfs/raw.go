package cbfs

import (
	"io"
	"log"
)

func init() {
	if err := RegisterFileReader(&SegReader{T: TypeRaw, N: "CBFSRaw", F: NewRaw}); err != nil {
		log.Fatal(err)
	}
}

func NewRaw(r CountingReader, f *File) (ReadWriter, error) {
	rec := &RawRecord{File: *f}
	Debug("Before Raw: total bytes read: %d", r.Count())
	Debug("Got header %v", *rec)
	return rec, nil
}

func (h *RawRecord) Read([]byte) (int, error) {
	return -1, nil
}

func (h *RawRecord) Write([]byte) (int, error) {
	return -1, nil
}

func (h *RawRecord) String() string {
	return recString(h.File.Name, h.RomOffset, h.Type.String(), h.Size, "none")
}

func (h *RawRecord) Name() string {
	return h.File.Name
}

func (r *RawRecord) Update(w io.Writer) error {
	if err := Write(w, r.FileHeader); err != nil {
		return err
	}
	return Write(w, r.Data)
}

func (r *RawRecord) Header() *File {
	return &r.File
}
