package cbfs

import (
	"fmt"
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

func (h *MasterRecord) Read([]byte) (int, error) {
	return -1, nil
}

func (h *MasterRecord) Write([]byte) (int, error) {
	return -1, nil
}

func (h *MasterRecord) String() string {
	return fmt.Sprintf("%s\t%#x\t%s\t%d\t%s", h.Name, h.SubHeaderOffset, h.Type.String(), h.Size, "none")
}
