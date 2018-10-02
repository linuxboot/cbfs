package cbfs

import (
	"fmt"
	"io"
	"log"
)

func init() {
	if err := RegisterFileReader(&SegReader{Type: TypeRaw, Name: "CBFSRaw", New: NewRaw}); err != nil {
		log.Fatal(err)
	}
}

func NewRaw(f *File) (ReadWriter, error) {
	rec := &RawRecord{File: *f}
	return rec, nil
}

func (r *RawRecord) Read(in io.ReadSeeker) error {
	_, err := in.Read(r.Data)
	if err != nil {
		return fmt.Errorf("raw read: %v", err)
	}
	Debug("raw data read OK")
	return nil
}

func (r *RawRecord) String() string {
	return recString(r.File.Name, r.RecordStart, r.Type.String(), r.Size, "none")
}

func (r *RawRecord) Write(w io.Writer) error {
	return Write(w, r.Data)
}

func (r *RawRecord) Header() *File {
	return &r.File
}
