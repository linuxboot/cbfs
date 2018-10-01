package cbfs

import (
	"io"
	"log"
)

func init() {
	if err := RegisterFileReader(&SegReader{Type: TypeCMOSLayout, Name: "CBFSCMOSLayout", New: NewCMOSLayout}); err != nil {
		log.Fatal(err)
	}
}

func NewCMOSLayout(f *File) (ReadWriter, error) {
	rec := &CMOSLayoutRecord{File: *f}
	Debug("Got header %v", *rec)
	return rec, nil
}

func (r *CMOSLayoutRecord) Read(in io.ReadSeeker) error {
	return nil
}

func (r *CMOSLayoutRecord) String() string {
	return recString(r.File.Name, r.RecordStart, r.Type.String(), r.Size, "none")
}

func (r *CMOSLayoutRecord) Write(w io.Writer) error {
	if err := Write(w, r.FileHeader); err != nil {
		return err
	}
	return Write(w, r.Data)
}

func (r *CMOSLayoutRecord) Header() *File {
	return &r.File
}
