package cbfs

import (
	"encoding/binary"
	"fmt"
	"io"
)

type SegReader struct {
	T CBFSFileType
	F func(r io.Reader) (CBFSReadWriter, error)
	N string
}

var SegReaders = make(map[CBFSFileType] *SegReader)

func RegisterFileReader(f *SegReader) error {
	if r, ok := SegReaders[f.T]; ok {
		return fmt.Errorf("RegisterFileType: Slot of %v is owned by %s, can't add %s", r.T, r.N, f.N)
	}
	SegReaders[f.T] = f
	Debug("Registered %v", f)
	return nil
}

func NewSegs(r io.ReaderAt) ([]CBFSReadWriter, error) {
	var off int64
	var segs []CBFSReadWriter
	for {
		var f File
		if err := binary.Read(io.NewSectionReader(r, off, TagSize), CBFSEndian, &f.LarchiveTag); err != nil {
			return nil, err
		}
		Debug("Found segment at %#x: %v", off, f.LarchiveTag)
		if string(f.LarchiveTag.Magic[:]) != FileMagic {
			Debug("It is not an LARCHIVE")
			off += 16
			continue
		}
		Debug("It is an LARCHIVE")
		if err := binary.Read(io.NewSectionReader(r, off, FileSize), CBFSEndian, &f); err != nil {
			return nil, err
		}
		Debug("It is type %v", f.Type)
		// If we match something, cons up a SectionReader for it and let the appropriate
		// type read it in.
		n, ok := SegReaders[f.Type]
		if ! ok {
			return nil, fmt.Errorf("%v: unknown type %v", f, f.Type)
		}
		Debug("Found a SegReader for this %d size section: %v", f.Size, n)
		s, err := n.F(io.NewSectionReader(r, off, int64(f.Size)))
		if err != nil {
			return nil, err
		}
		Debug("Segment was readable")
		segs = append(segs, s)
	}
	return segs, nil
}
