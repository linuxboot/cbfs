package cbfs

import (
	"fmt"
	"io"
	"log"
)

func init() {
	if err := RegisterFileReader(&SegReader{Type: TypeStage, Name: "Stage", New: NewStageRecord}); err != nil {
		log.Fatal(err)
	}
}

func NewStageRecord(f *File) (ReadWriter, error) {
	r := &StageRecord{File: *f}
	return r, nil
}

func (r *StageRecord) Read(in io.ReadSeeker) error {
	if err := ReadLE(in, &r.StageHeader); err != nil {
		Debug("StageHeader read: %v", err)
		return err
	}
	Debug("Got StageHeader %s, data is %d bytes", r.String(), r.StageHeader.Size)
	r.Data = make([]byte, r.StageHeader.Size)
	n, err := in.Read(r.Data)
	if err != nil {
		return err
	}
	Debug("Stage read %d bytes", n)
	return nil
}

func (h *StageHeader) String() string {
	return fmt.Sprintf("Compression %#x Entry %#x LoadAddress %#x Size %#x MemSize %#x",
		h.Compression,
		h.Entry,
		h.LoadAddress,
		h.Size,
		h.MemSize)
}

func (h *StageRecord) String() string {
	return recString(h.File.Name, h.RecordStart, h.Type.String(), h.Size, h.Compression.String())
}

func (r *StageRecord) Write(w io.Writer) error {
	if err := WriteLE(w, r.StageHeader); err != nil {
		return err
	}

	return Write(w, r.Data)
}

func (r *StageRecord) File() *File {
	return &r.File
}
