package cbfs

import (
	"fmt"
	"io"
)

type SegReader struct {
	T CBFSFileType
	F func(r io.Reader, f *File) (CBFSReadWriter, error)
	N string
}

var SegReaders = make(map[CBFSFileType]*SegReader)

func RegisterFileReader(f *SegReader) error {
	if r, ok := SegReaders[f.T]; ok {
		return fmt.Errorf("RegisterFileType: Slot of %v is owned by %s, can't add %s", r.T, r.N, f.N)
	}
	SegReaders[f.T] = f
	Debug("Registered %v", f)
	return nil
}

func NewSegs(r io.Reader) ([]CBFSReadWriter, error) {
	var segs []CBFSReadWriter
	for {
		var f File
		var m Magic
		err := CBFSRead(r, m[:])
		if err == io.EOF {
			return segs, nil
		}
		if err != nil {
			return nil, err
		}
		if string(m[:]) != FileMagic {
			// Do a fake read and throw away the results.
			err := CBFSRead(r, m[:])
			if err == io.EOF {
				return segs, nil
			}
			if err != nil {
				return nil, err
			}

			continue
		}
		Debug("It is an LARCHIVE")
		if err := CBFSRead(r, &f); err != nil {
			return nil, err
		}
		Debug("It is %v type %v", f, f.Type)
		// If we match something, cons up a SectionReader for it and let the appropriate
		// type read it in.
		n, ok := SegReaders[f.Type]
		if !ok {
			return nil, fmt.Errorf("%v: unknown type %v", f, f.Type)
		}
		Debug("Found a SegReader for this %d size section: %v", f.Size, n)
		s, err := n.F(r, &f)
		if err != nil {
			return nil, err
		}
		Debug("Segment was readable")
		segs = append(segs, s)
	}
	return segs, nil
}
