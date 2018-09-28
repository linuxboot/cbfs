package cbfs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

var Debug = func(format string, v ...interface{}) {}

// Read reads things in in BE format, which they are supposed to be in.
func Read(r io.Reader, f interface{}) error {
	if err := binary.Read(r, Endian, f); err != nil {
		return err
	}
	return nil
}

// ReadLE reads things in LE format, which the spec says it is not in.
func ReadLE(r io.Reader, f interface{}) error {
	if err := binary.Read(r, binary.LittleEndian, f); err != nil {
		return err
	}
	return nil
}

// Write reads things in in BE format, which they are supposed to be in.
func Write(w io.Writer, f interface{}) error {
	if err := binary.Writer(r, Endian, f); err != nil {
		return err
	}
	return nil
}

// WriteLE reads things in LE format, which the spec says it is not in.
func WriteLE(r io.Writer, f interface{}) error {
	if err := binary.Writer(r, binary.LittleEndian, f); err != nil {
		return err
	}
	return nil
}

func ReadName(r CountingReader, f *File, size uint32) (string, error) {
	b := make([]byte, size)
	n, err := r.Read(b)
	if err != nil {
		Debug("ReadName failed:%v", err)
		return "", err
	}
	Debug("Readname gets %#02x", b)
	if n != len(b) {
		err = fmt.Errorf("ReadName: got %d, want %d for name", n, len(b))
		Debug("Readname short: %v", err)
		return "", err
	}
	// discard trailing NULLs
	z := bytes.Split(b, []byte{0})
	return string(z[0]), nil
}

func Align(r CountingReader) error {
	var junk [16]byte
	align := (int(r.Count()) + 15) & ^0xf
	amt := align - int(r.Count())
	return Read(r, junk[:amt])
}

func (c Compression) String() string {
	switch c {
	case None:
		return "none"
	case LZMA:
		return "lzma"
	case LZ4:
		return "lz4"
	}
	return "unknown"
}

func (f FileType) String() string {
	switch f {
	case TypeDeleted2:
		return "TypeDeleted2"
	case TypeDeleted:
		return "TypeDeleted"
	case TypeMaster:
		return "cbfs header"
	case TypeBootBlock:
		return "TypeBootBlock"
	case TypeStage:
		return "TypeStage"
	case TypeSELF:
		return "TypeSELF"
	case TypeFIT:
		return "TypeFIT"
	case TypeOptionRom:
		return "TypeOptionRom"
	case TypeBootSplash:
		return "TypeBootSplash"
	case TypeRaw:
		return "TypeRaw"
	case TypeVSA:
		return "TypeVSA"
	case TypeMBI:
		return "TypeMBI"
	case TypeMicroCode:
		return "TypeMicroCode"
	case TypeFSP:
		return "TypeFSP"
	case TypeMRC:
		return "TypeMRC"
	case TypeMMA:
		return "TypeMMA"
	case TypeEFI:
		return "TypeEFI"
	case TypeStruct:
		return "TypeStruct"
	case TypeCMOS:
		return "TypeCMOS"
	case TypeSPD:
		return "TypeSPD"
	case TypeMRCCache:
		return "TypeMRCCache"
	case TypeCMOSLayout:
		return "TypeCMOSLayout"
	}
	return fmt.Sprintf("%#x", uint32(f))
}

func recString(n string, off uint32, typ string, sz uint32, compress string) string {
	return fmt.Sprintf("%s\t\t%#x\t%s\t%d\t%s", n, off, typ, sz, compress)
}
