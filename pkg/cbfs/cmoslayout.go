package cbfs

import (
	"io"
	"log"
)

func init() {
	if err := RegisterFileReader(&SegReader{T: TypeCMOSLayout, N: "CBFSCMOSLayout", F: NewCMOSLayout}); err != nil {
		log.Fatal(err)
	}
}

func NewCMOSLayout(r CountingReader, f *File) (ReadWriter, error) {
	rec := &CMOSLayoutRecord{File: *f}
	Debug("Before CMOSLayout: total bytes read: %d", r.Count())
	Debug("Got header %v", *rec)
	return rec, nil
}

func (r *CMOSLayoutRecord) Read([]byte) (int, error) {
	return -1, nil
}

func (r *CMOSLayoutRecord) Write([]byte) (int, error) {
	return -1, nil
}

func (r *CMOSLayoutRecord) String() string {
	return recString(r.File.Name, r.RomOffset, r.Type.String(), r.Size, "none")
}

func (r *CMOSLayoutRecord) Name() string {
	return r.File.Name
}

func (r *CMOSLayoutRecord) Update(w io.Writer) error {
	if err := Write(w, r.FileHeader); err != nil {
		return err
	}
	return Write(w, r.Data)
}

func (r *CMOSLayoutRecord) Header() *File {
	return &r.File
}
