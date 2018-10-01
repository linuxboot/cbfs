package cbfs

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/linuxboot/fiano/pkg/fmap"
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

func NewImage(in io.ReadSeeker) (*Image, error) {
	// Suck the image in. Todo: write a thing that implements
	// ReadSeeker on a []byte.
	b, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	if _, err := in.Seek(0, 0); err != nil {
		return nil, err
	}
	f, m, err := fmap.Read(in)
	if err != nil {
		return nil, err
	}
	Debug("Fmap %v", f)
	var i = &Image{Offset: -1, FMAP: f, FMAPMetadata: m, Data: b}
	var x = int(-1)
	for i, a := range f.Areas {
		Debug("Check %v", a.Name.String())
		if a.Name.String() == "COREBOOT" {
			x = i
			break
		}
	}
	if x == -1 {
		return nil, fmt.Errorf("No CBFS in fmap")
	}
	i.Index = x
	Debug("COREBOOT is the %d entry", x)
	fr, err := f.ReadArea(in, x)
	if err != nil {
		return nil, err
	}
	r := NewCountingReader(fr)

	for {
		var f File
		var m Magic
		err := Align(r)
		if err == io.EOF {
			return i, nil
		}
		if err != nil {
			return nil, err
		}
		recStart := r.Count()
		err = Read(r, m[:])
		if err == io.EOF {
			return i, nil
		}
		if err != nil {
			return nil, err
		}
		if string(m[:]) != FileMagic {
			continue
		}
		Debug("It is an LARCHIVE at %#x", int(r.Count())-len(FileMagic))
		if i.Offset < 0 {
			i.Offset = int(recStart)
		}
		if err := Read(r, &f.FileHeader); err != nil {
			Debug("Reading the File failed: %v", err)
			return nil, err
		}
		f.RomOffset = recStart
		Debug("It is %v type %v", f, f.Type)
		sr, ok := SegReaders[f.Type]
		if !ok {
			return nil, fmt.Errorf("%v: unknown type %v", f, f.Type)
		}
		headSize := r.Count() - recStart
		Debug("Namelen %d %d ", f.SubHeaderOffset, f.SubHeaderOffset-headSize)
		n, err := ReadName(r, &f, f.SubHeaderOffset-headSize)
		if err != nil {
			return nil, err
		}
		f.Name = n
		Debug("Count after name is %#x", r.Count())
		Debug("Found a SegReader for this %d size section: %v", f.Size, n)
		s, err := sr.F(r, &f)
		if err != nil {
			return nil, err
		}
		Debug("Segment was readable")
		i.Segs = append(i.Segs, s)
		Debug("r.Count is now %#x", r.Count())
	}
	return i, nil
}

func (i *Image) WriteFile(name string, perm os.FileMode) error {
	if err := ioutil.WriteFile(name, i.Data, 0666); err != nil {
		return err
	}
	return nil
}

// Update creates a new []byte for the cbfs. It is complicated a lot
// by the fact that endianness is not consistent in cbfs images.
func (i *Image) Update() error {
	// Because there can be gaps due to alignment of various
	// components, we start out by filling i.Data with ff
	// past the FMAP header.
	for x := range(i.Data[512:]) {
		i.Data[x+512] = 0xff
	}
	for _, s := range i.Segs {
		var b bytes.Buffer
		if err := s.Update(&b); err != nil {
			return err
		}
		Debug("Copy %d bytes to i.Data[%d]", len(b.Bytes()), s.Header().RomOffset+512)
		copy(i.Data[s.Header().RomOffset+512:], b.Bytes())
	}
	return nil
}

func (i *Image) String() string {
	var s = "FMAP REGION: COREBOOT\nName\t\t\t\tOffset\tType\t\tSize\tComp\n"
	for _, seg := range i.Segs {
		s = s + seg.String() + "\n"
	}
	return s
}

func (i *Image) Remove(n string) error {
	found := -1
	for x, s := range i.Segs {
		if s.Name() == n {
			found = x
		}
	}
	if found == -1 {
		return os.ErrExist
	}
	// You can not remove the master header or the boot block.
	// Just remake the cbfs if you're doing that kind of surgery.
	if found == 0 || found == len(i.Segs)-1 {
		return os.ErrPermission
	}
	del := &EmptyRecord{File: *i.Segs[found].Header(), Data: make([]byte, i.Segs[found].Header().Size)}
	del.Type = TypeDeleted2
	i.Segs[found] = del
	/* not yet
	// We might be able to merge it. This is not common however.
	if i.Segs[found].Type != i.Segs[found+1].Type {
		return nil
	}
	end := i.Segs[found].Offset + i.Segs[found].Size
	if i.Segs[found+1].Offset != end {
		return nil
	}
	i.Segs[found+1].Offset = i.Segs[found].Offset
	i.Segs[found+1].Size += i.Segs[found].Size
	i.Segs = append(i.Segs[:found], i.Segs[found:]...)
	*/
	return nil
}
