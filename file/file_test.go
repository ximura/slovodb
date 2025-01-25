package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileMgrWriteRead(t *testing.T) {
	blockSize := 16
	dataDir := "testdata"
	testFile := "writetestfile"

	mgr := NewFileMgr(dataDir, blockSize)
	t.Cleanup(func() {
		mgr.Close()
		os.Remove(filepath.Join(dataDir, testFile))
	})

	// write initial block
	blockZero := &BlockID{
		Filename: testFile,
		Number:   0,
	}
	data := "aaaaaaaaaaaaaaaa"
	checkWrite(t, mgr, blockZero, data)
	checkRead(t, mgr, blockZero, data)
	checkFileContent(t, filepath.Join(dataDir, testFile), "aaaaaaaaaaaaaaaa")

	// Write to block 1
	blockOne := &BlockID{
		Filename: testFile,
		Number:   1,
	}
	data = "bbbbbbbbbbbbbbbb"
	checkWrite(t, mgr, blockOne, data)
	checkRead(t, mgr, blockOne, data)
	checkFileContent(t, filepath.Join(dataDir, testFile), "aaaaaaaaaaaaaaaabbbbbbbbbbbbbbbb")

	// Rewrite to block 0
	data = "cccccccccccccccc"
	checkWrite(t, mgr, blockZero, data)
	checkRead(t, mgr, blockZero, data)
	checkFileContent(t, filepath.Join(dataDir, testFile), "ccccccccccccccccbbbbbbbbbbbbbbbb")
}

func checkWrite(t *testing.T, mgr *FileMgr, blockID *BlockID, data string) {
	page := NewPage(mgr.blockSize)
	page.Write(0, []byte(data))

	n, err := mgr.Write(blockID, page)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != mgr.blockSize {
		t.Fatalf("Write returned %d, want %d", n, mgr.blockSize)
	}
}

func checkRead(t *testing.T, mgr *FileMgr, blockID *BlockID, want string) {
	page := NewPage(mgr.blockSize)
	n, err := mgr.Read(blockID, page)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if n != mgr.blockSize {
		t.Fatalf("Read returned %d, want %d", n, mgr.blockSize)
	}
	if string(page.Bytes()) != want {
		t.Fatalf("Read returned %q, want %q", page.Bytes(), want)
	}
}

func checkFileContent(t *testing.T, filename, want string) {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	got := make([]byte, len(want))
	_, err = f.Read([]byte(got))
	if err != nil && err.Error() != "EOF" {
		t.Fatalf("Failed to read file: %v", err)
	}
	defer f.Close()

	if string(got) != want {
		t.Fatalf("File content is %q, want %q", got, want)
	}
}
