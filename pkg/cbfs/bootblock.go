package cbfs

import (
	"io"
	"log"
)

func init() {
	if err := RegisterFileReader(&SegReader{T: TypeBootBlock, N: "CBFSBootBlock", F: NewBootBlock}); err != nil {
		log.Fatal(err)
	}
}

func NewBootBlock(r CountingReader, f *File) (ReadWriter, error) {
	h := &BootBlockRecord{File: *f}
	Debug("Before BootBlock: total bytes read: %d", r.Count())
	Debug("Got header %v", *h)
	h.Data = make([]byte, h.Size)
	n, err := r.Read(h.Data)
	if err != nil {
		return nil, err
	}
	Debug("Bootblock read %d bytes", n)
	return h, nil
}

func (r *BootBlockRecord) Read([]byte) (int, error) {
	return -1, nil
}

func (r *BootBlockRecord) Write([]byte) (int, error) {
	return -1, nil
}

func (r *BootBlockRecord) String() string {
	return recString(r.Name(), r.RomOffset, r.Type.String(), r.Size, "none")
}

func (r *BootBlockRecord) Name() string {
	return "BootBlock"
}

func (r *BootBlockRecord) Update(w io.Writer) error {
	if err := Write(w, r.FileHeader); err != nil {
		return err
	}
	return Write(w, r.Data)
}

func (r *BootBlockRecord) Header() *File {
	return &r.File
}
