package cbfs

import (
	"io"
	"log"
)

func init() {
	if err := RegisterFileReader(&SegReader{T: 2, N: "CBFSMaster", F: NewMaster}); err != nil {
		log.Fatal(err)
	}
}

func NewMaster(r CountingReader, f *File) (ReadWriter, error) {
	h := &MasterRecord{File: *f}
	Debug("Before Header: total bytes read: %d", r.Count())
	if err := Read(r, &h.MasterHeader); err != nil {
		Debug("Header read: %v", err)
		return nil, err
	}
	Debug("Got header %s offset %#x", h.String(), h.Offset)
	return h, nil
}

func (r *MasterRecord) Read([]byte) (int, error) {
	return -1, nil
}

func (r *MasterRecord) Write([]byte) (int, error) {
	return -1, nil
}

func (r *MasterRecord) String() string {
	return recString(r.File.Name, r.RomOffset, r.Type.String(), r.Size, "none")
}

func (r *MasterRecord) Name() string {
	return r.File.Name
}

func (r *MasterRecord) Update(w io.Writer) error {
	return Write(w, r.MasterHeader)
}

func (r *MasterRecord) Header() *File {
	return &r.File
}
