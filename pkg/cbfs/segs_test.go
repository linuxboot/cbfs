package cbfs

import (
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	f, err := os.Open("testdata/coreboot.rom")
	if err != nil {
		t.Fatal(err)
	}
	s, err := NewSegs(f)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Segs %v", s)
}
