package cbfs

import (
	"fmt"
	"io"
)

type SegReader struct {
	T FileType
	F func(r CountingReader, f *File) (ReadWriter, error)
	N string
}

var SegReaders = make(map[FileType]*SegReader)

func RegisterFileReader(f *SegReader) error {
	if r, ok := SegReaders[f.T]; ok {
		return fmt.Errorf("RegisterFileType: Slot of %v is owned by %s, can't add %s", r.T, r.N, f.N)
	}
	SegReaders[f.T] = f
	Debug("Registered %v", f)
	return nil
}

func NewSegs(in io.Reader) ([]ReadWriter, error) {
	r := NewCountingReader(in)
	var segs []ReadWriter
	for {
		var f File
		var m Magic
		err := Read(r, m[:])
		if err == io.EOF {
			return segs, nil
		}
		if err != nil {
			return nil, err
		}
		if string(m[:]) != FileMagic {
			// Do a fake read and throw away the results.
			err := Read(r, m[:])
			if err == io.EOF {
				return segs, nil
			}
			if err != nil {
				return nil, err
			}

			continue
		}
		Debug("It is an LARCHIVE")
		if err := Read(r, &f.FileHeader); err != nil {
			Debug("Reading the File failed: %v", err)
			return nil, err
		}
		Debug("It is %v type %v", f, f.Type)
		Debug("%d %d ", f.Offset, f.Offset - 24)
		n, err := ReadName(r, &f)
		if err != nil {
			return nil, err
		}
		f.Name = n

		sr, ok := SegReaders[f.Type]
		if !ok {
			return nil, fmt.Errorf("%v: unknown type %v", f, f.Type)
		}
		Debug("Found a SegReader for this %d size section: %v", f.Size, n)
		s, err := sr.F(r, &f)
		if err != nil {
			return nil, err
		}
		Debug("Segment was readable")
		segs = append(segs, s)
	}
	return segs, nil
}
