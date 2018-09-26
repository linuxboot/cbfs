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
		b    []byte
		want string
	}{
		{"Master Only", Master, ""},
	}

	for _, tc := range tests {
		t.Run(tc.n, func(t *testing.T) {
			r := bytes.NewReader(tc.b)
			_, err := NewSegs(r)
			if err != nil {
				t.Errorf("got %v, want nil", err)
				return
			}
		})
	}
}

func TestConflict(t *testing.T) {
	if err := RegisterFileReader(&SegReader{T: 2, N: "CBFSRaw", F: nil}); err == nil {
		t.Fatalf("Registering conflicting entry to type 2, want error, got nil")
	}

}
