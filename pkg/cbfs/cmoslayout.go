package cbfs

import (
	"fmt"
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
	_, err := in.Read(r.Data)
	if err != nil {
		return fmt.Errorf("CMOSLayout: empty read: %v", err)
	}
	Debug("CMOSLayout data read OK")
	return nil

	return nil
}

func (r *CMOSLayoutRecord) String() string {
	return recString(r.File.Name, r.RecordStart, r.Type.String(), r.Size, "none")
}

func (r *CMOSLayoutRecord) Write(w io.Writer) error {
	return Write(w, r.Data)
}

func (r *CMOSLayoutRecord) Header() *File {
	return &r.File
}
