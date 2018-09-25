package cbfs

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {
	f, err := os.Open("testdata/coreboot.rom")
	if err != nil {
		t.Fatal(err)
	}
	s, err := NewSegs(f)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", s)
}

func TestBogusArchives(t *testing.T) {
	var tests = []struct {
		n    string
		r    io.Reader
		want string
	}{
		{"Short", bytes.NewReader([]byte("INUXARCHIV")), "unexpected EOF"},
		{"Misaligned", bytes.NewReader([]byte("INUXARCHIVL")), "unexpected EOF"},
	}

	for _, tc := range tests {
		t.Run(tc.n, func(t *testing.T) {
			_, err := NewSegs(tc.r)
			if err == nil {
				t.Errorf("got nil, want %v", tc.want)
				return
			}
			e := fmt.Sprintf("%v", err)
			if e != tc.want {
				t.Errorf("got %v, want %v", e, tc.want)
			}
		})
	}
}

func TestReadSimple(t *testing.T) {
	var tests = []struct {
		n    string
		b []byte
		want string
	}{
		{"Master Only", []byte{0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x38, 0x63, 0x62, 0x66, 0x73, 0x20, 0x06d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x20, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ""},
	}

	for _, tc := range tests {
		t.Run(tc.n, func(t *testing.T) {
			r := bytes.NewReader(append([]byte(FileMagic), tc.b...))
			_, err := NewSegs(r)
			if err != nil {
				t.Errorf("got %v, want nil", err)
				return
			}
		})
	}
}
