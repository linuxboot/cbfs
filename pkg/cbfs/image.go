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
	Type FileType
	New  func(f *File) (ReadWriter, error)
	Name string
}

var SegReaders = make(map[FileType]*SegReader)

func RegisterFileReader(f *SegReader) error {
	if r, ok := SegReaders[f.Type]; ok {
		return fmt.Errorf("RegisterFileType: Slot of %v is owned by %s, can't add %s", r.Type, r.Name, f.Name)
	}
	SegReaders[f.Type] = f
	Debug("Registered %v", f)
	return nil
}

func NewImage(rs io.ReadSeeker) (*Image, error) {
	// Suck the image in. Todo: write a thing that implements
	// ReadSeeker on a []byte.
	b, err := ioutil.ReadAll(rs)
	if err != nil {
		return nil, fmt.Errorf("ReadAll: %v", err)
	}
	in := bytes.NewReader(b)
	f, m, err := fmap.Read(in)
	if err != nil {
		return nil, err
	}
	Debug("Fmap %v", f)
	var i = &Image{FMAP: f, FMAPMetadata: m, Data: b}
	for _, a := range f.Areas {
		Debug("Check %v", a.Name.String())
		if a.Name.String() == "COREBOOT" {
			i.Area = &a
			break
		}
	}
	if i.Area == nil {
		return nil, fmt.Errorf("No CBFS in fmap")
	}
	r := io.NewSectionReader(in, int64(i.Area.Offset), int64(i.Area.Size))

	for off := int64(0); off < int64(i.Area.Size); {
		var f File
		if _, err := r.Seek(off, io.SeekStart); err != nil {
			return nil, err
		}
		err := Read(r, &f.FileHeader)
		if err == io.EOF {
			return i, nil
		}
		if err != nil {
			return nil, err
		}
		if string(f.Magic[:]) != FileMagic {
			off += 16
			continue
		}
		Debug("It is %v type %v", f, f.Type)
		f.RecordStart = uint32(off)
		nameStart, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("Getting file offset for name: %v", err)
		}
		sr, ok := SegReaders[f.Type]
		if !ok {
			return nil, fmt.Errorf("%v: unknown type %v", f, f.Type)
		}
		n, err := ReadName(r, &f, f.SubHeaderOffset-(uint32(nameStart)-f.RecordStart))
		if err != nil {
			return nil, err
		}
		f.Name = n
		Debug("Found a SegReader for this %d size section: %v", f.Size, n)
		s, err := sr.New(&f)
		if err != nil {
			return nil, err
		}
		if err := s.Read(r); err != nil {
			return nil, err
		}
		Debug("Segment was readable")
		i.Segs = append(i.Segs, s)
		off, err = r.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, err
		}
		// Force alignment.
		off = (off + 15) & (^15)

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
	for x := range i.Data[512:] {
		i.Data[x+512] = 0xff
	}
	for _, s := range i.Segs {
		var b bytes.Buffer
		if err := Write(&b, s.Header().FileHeader); err != nil {
			return err
		}
		// Tradition seems to have it that name bytes are zero-filled, not
		// 0xff-filled. That's stupid.
		n := make([]byte, s.Header().FileHeader.SubHeaderOffset-FileSize)
		for i := range n {
			n[i] = 0xff
		}
		n[len(s.Header().Name)] = 0
		copy(n, []byte(s.Header().Name))
		if _, err := b.Write(n); err != nil {
			return fmt.Errorf("Writing name to cbfs record for %v: %v", s, err)
		}
		if err := s.Write(&b); err != nil {
			return err
		}
		Debug("Copy %d bytes to i.Data[%d]", len(b.Bytes()), s.Header().RecordStart+512)
		copy(i.Data[s.Header().RecordStart+512:], b.Bytes())
	}
	return nil
}

func (i *Image) String() string {
	var s = "FMAP REGIOName: COREBOOT\nName\t\t\t\tOffset\tType\t\tSize\tComp\n"
	for _, seg := range i.Segs {
		s = s + seg.String() + "\n"
	}
	return s
}

func (i *Image) Remove(n string) error {
	found := -1
	for x, s := range i.Segs {
		if s.Header().Name == n {
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
	del, _ := NewEmptyRecord(i.Segs[found].Header())
	i.Segs[found] = del
	return nil
}
