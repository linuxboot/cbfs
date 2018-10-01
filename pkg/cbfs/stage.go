package cbfs

import (
	"fmt"
	"io"
	"log"
)

func init() {
	if err := RegisterFileReader(&SegReader{T: TypeStage, N: "Stage", F: NewStageRecord}); err != nil {
		log.Fatal(err)
	}
}

func NewStageRecord(r CountingReader, f *File) (ReadWriter, error) {
	h := &StageRecord{File: *f}
	Debug("Before StageRecord: total bytes read: %d", r.Count())
	if err := ReadLE(r, &h.StageHeader); err != nil {
		Debug("StageHeader read: %v", err)
		return nil, err
	}
	Debug("Got StageHeader %s, data is %d bytes", h.String(), h.StageHeader.Size)
	h.Data = make([]byte, h.StageHeader.Size)
	n, err := r.Read(h.Data)
	if err != nil {
		return nil, err
	}
	Debug("Stage read %d bytes", n)
	return h, nil
}

func (h *StageRecord) Read([]byte) (int, error) {
	return -1, nil
}

func (h *StageRecord) Write([]byte) (int, error) {
	return -1, nil
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
	return recString(h.File.Name, h.RomOffset, h.Type.String(), h.Size, h.Compression.String())
}

func (h *StageRecord) Name() string {
	return h.File.Name
}

func (r *StageRecord) Update(w io.Writer) error {
	if err := Write(w, r.FileHeader); err != nil {
		return err
	}
	if err := WriteLE(w, r.StageHeader); err != nil {
		return err
	}

	return Write(w, r.Data)
}

func (r *StageRecord) Header() *File {
	return &r.File
}
