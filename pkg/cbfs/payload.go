package cbfs

import (
	"fmt"
	"log"
)

func init() {
	if err := RegisterFileReader(&SegReader{T: TypeSELF, N: "Payload", F: NewPayloadRecord}); err != nil {
		log.Fatal(err)
	}
}

func NewPayloadRecord(r CountingReader, f *File) (ReadWriter, error) {
	p := &PayloadRecord{File: *f}
	Debug("Before PayloadRecord: total bytes read: %d", r.Count())
	for {
		var h PayloadHeader
		if err := Read(r, &h); err != nil {
			Debug("PayloadHeader read: %v", err)
			return nil, err
		}
		Debug("Got PayloadHeader %s", h.String())
		p.Segs = append(p.Segs, h)
		if h.Type == SegEntry {
			break
		}
	}
	p.Data = make([]byte, p.Size)
	n, err := r.Read(p.Data)
	if err != nil {
		return nil, err
	}
	Debug("Payload read %d bytes", n)
	return p, nil
}

func (h *PayloadRecord) Read([]byte) (int, error) {
	return -1, nil
}

func (h *PayloadRecord) Write([]byte) (int, error) {
	return -1, nil
}

func (h *PayloadRecord) String() string {
	s := recString(h.Name, h.RomOffset, h.Type.String(), h.Size, "none")
	for i, seg := range h.Segs {
		s += recString(fmt.Sprintf("\n\tSeg #%d\t", i), seg.Offset, "Payload segment", seg.Size, seg.Compression.String())
	}
	return s
}

func (h *PayloadHeader) String() string {
	return fmt.Sprintf("Type %#x Compression %#x Offset %#x LoadAddress %#x Size %#x MemSize %#x",
		h.Type,
		h.Compression,
		h.Offset,
		h.LoadAddress,
		h.Size,
		h.MemSize)
}
