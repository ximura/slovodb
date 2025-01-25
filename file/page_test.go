package file

import (
	"errors"
	"testing"
)

func TestPageWriteREad(t *testing.T) {
	page := NewPage(16)

	// Write data to the page
	data := []byte("Hello, world")
	n, err := page.Write(0, data)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != len(data) {
		t.Fatalf("Write returned %d, expected %d", n, len(data))
	}

	// Read the page from the beginning
	got := make([]byte, len(data))
	page.Read(0, got)
	if string(got) != string(data) {
		t.Fatalf("Read failed: got %q, want %q", string(got), string(data))
	}

	//Write data at an offset
	data = []byte("SlovoDB!!")
	n, err = page.Write(7, data)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != len(data) {
		t.Fatalf("Write returned %d, expected %d", n, len(data))
	}

	// Read the entire page
	want := []byte("Hello, SlovoDB!!")
	got = page.Bytes()
	if string(got) != string(want) {
		t.Fatalf("Read failed: got %q, want %q", string(got), string(want))
	}

	// Read only a section of the page
	got = make([]byte, len(data))
	page.Read(7, got)
	if string(got) != string(data) {
		t.Fatalf("Read failed: got %q, want %q", string(got), string(want))
	}

	// Write data bigger than page size
	data = []byte("longer data")
	_, err = page.Write(10, data)
	expectedErr := errors.New("data exceeds page bounds")
	if err == nil || err.Error() != expectedErr.Error() {
		t.Fatalf("Write returrned no error or what was not expected: got %q, expected %q", err, expectedErr)
	}
}
