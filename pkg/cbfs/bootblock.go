package cbfs

import (
	"fmt"
	"io"
	"log"
)

func init() {
	if err := RegisterFileReader(&SegReader{Type: TypeBootBlock, Name: "CBFSBootBlock", New: NewBootBlock}); err != nil {
		log.Fatal(err)
	}
}

func NewBootBlock(f *File) (ReadWriter, error) {
	r := &BootBlockRecord{File: *f}
	Debug("Got header %v", *r)
	r.Data = make([]byte, r.Size)
	return r, nil
}

func (r *BootBlockRecord) Read(in io.ReadSeeker) error {
	n, err := in.Read(r.Data)
	if err != nil {
		return fmt.Errorf("Reading bootblockrecord at %#x, got %d bytes, wanted %d", r.SubHeaderOffset, n, len(r.Data))
	}
	Debug("Bootblock read %d bytes", n)
	return nil
}

func (r *BootBlockRecord) String() string {
	return recString(r.File.Name, r.RecordStart, r.Type.String(), r.Size, "none")
}

func (r *BootBlockRecord) Write(w io.Writer) error {
	return Write(w, r.Data)
}

func (r *BootBlockRecord) Header() *File {
	return &r.File
}
