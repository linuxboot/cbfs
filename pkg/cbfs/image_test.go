package cbfs

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestReadFile(t *testing.T) {
	Debug = t.Logf
	f, err := os.Open("testdata/coreboot.rom")
	if err != nil {
		t.Fatal(err)
	}
	i, err := NewImage(f)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", i)
}

func TestBogusArchives(t *testing.T) {
	var tests = []struct {
		n    string
		r    io.ReadSeeker
		want string
	}{
		{"Short", bytes.NewReader([]byte("INUXARCHIV")), "Cannot find fmap signature"},
		{"Misaligned", bytes.NewReader([]byte("INUXARCHIVL")), "Cannot find fmap signature"},
	}

	for _, tc := range tests {
		t.Run(tc.n, func(t *testing.T) {
			_, err := NewImage(tc.r)
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
	Debug = t.Logf
	for _, tc := range tests {
		t.Run(tc.n, func(t *testing.T) {
			r := bytes.NewReader(tc.b)
			_, err := NewImage(r)
			if err != nil {
				t.Errorf("got %v, want nil", err)
				return
			}
		})
	}
}

func TestConflict(t *testing.T) {
	if err := RegisterFileReader(&SegReader{Type: 2, Name: "CBFSRaw", New: nil}); err == nil {
		t.Fatalf("Registering conflicting entry to type 2, want error, got nil")
	}

}

func TestStringer(t *testing.T) {
	f, err := os.Open("testdata/coreboot.rom")
	if err != nil {
		t.Fatal(err)
	}
	i, err := NewImage(f)
	if err != nil {
		t.Fatal(err)
	}
	s := i.String()

	t.Logf("Image string: %v", s)
}

func TestSimpleWrite(t *testing.T) {
	Debug = t.Logf
	f, err := os.Open("testdata/coreboot.rom")
	if err != nil {
		t.Fatal(err)
	}
	i, err := NewImage(f)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	out, err := ioutil.TempFile("", "cbfs")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(out.Name())
	if err := i.WriteFile(out.Name(), 0666); err != nil {
		t.Fatal(err)
	}
	out.Close()

	fi, err := os.Stat(out.Name())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("output file %v", fi)
	old, err := ioutil.ReadFile("testdata/coreboot.rom")
	if err != nil {
		t.Fatal(err)
	}
	new, err := ioutil.ReadFile(out.Name())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(old, new) {
		t.Fatalf("testdata/coreboot.rom and %s differ", out.Name())
	}

}

func TestRemovePayload(t *testing.T) {
	Debug = t.Logf
	f, err := os.Open("testdata/coreboot.rom")
	if err != nil {
		t.Fatal(err)
	}
	i, err := NewImage(f)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	out, err := ioutil.TempFile("", "cbfs")
	if err != nil {
		t.Fatal(err)
	}
	//	defer os.Remove(out.Name())
	if err := i.Remove("fallback/payload"); err != nil {
		t.Fatal(err)
	}
	if err := i.Update(); err != nil {
		t.Fatal(err)
	}
	if err := i.WriteFile(out.Name(), 0666); err != nil {
		t.Fatal(err)
	}
	out.Close()

	fi, err := os.Stat(out.Name())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("output file %v", fi)
	f, err = os.Open("testdata/removepayload.rom")
	if err != nil {
		t.Fatal(err)
	}
	old, err := NewImage(f)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	f, err = os.Open(out.Name())
	if err != nil {
		t.Fatalf("%s: %v", out.Name(), err)
	}
	new, err := NewImage(f)
	if err != nil {
		t.Fatalf("%s: %v", out.Name(), err)
	}
	f.Close()
	if !reflect.DeepEqual(old, new) {
		t.Errorf("testdata/coreboot.rom and %s differ", out.Name())
	}
	t.Logf("new image is %s", new.String())

}
