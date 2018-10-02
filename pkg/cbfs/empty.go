package cbfs

import (
	"io"
	"log"
)

func init() {
	if err := RegisterFileReader(&SegReader{Type: TypeDeleted, Name: "CBFSEmpty", New: NewEmptyRecord}); err != nil {
		log.Fatal(err)
	}
	if err := RegisterFileReader(&SegReader{Type: TypeDeleted2, Name: "CBFSEmpty", New: NewEmptyRecord}); err != nil {
		log.Fatal(err)
	}
}

func NewEmptyRecord(f *File) (ReadWriter, error) {
	r := &EmptyRecord{File: *f}
	Debug("Got header %v", r.String())
	// A common way to create a new empty record is to delete a file.
	// For the case that this is a remove, i.e. the file type
	// is changing, we just set the type and that's it. That way
	// we avoid spurious flash write cycles.
	if f.Type != TypeDeleted2 && f.Type != TypeDeleted {
		f.Type = TypeDeleted2
		return r, nil
	}
	r.Type = TypeDeleted2
	r.Name = ""
	return r, nil
}

func (r *EmptyRecord) Read(in io.ReadSeeker) error {
	return nil
}

func (r *EmptyRecord) String() string {
	return recString("(empty)", r.RecordStart, r.Type.String(), r.Size, "none")
}

func (r *EmptyRecord) Write(w io.Writer) error {
	return Write(w, r.FData)
}

func (r *EmptyRecord) Header() *File {
	return &r.File
}
