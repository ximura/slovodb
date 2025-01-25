package file

import "errors"

type Page struct {
	bytes []byte
}

// NewPage creates a new page with specified size
func NewPage(size int) *Page {
	return &Page{
		make([]byte, size),
	}
}

// Write copies data from data slice into page with specified offset
func (p *Page) Write(offset int, data []byte) (int, error) {
	if offset+len(data) > p.Size() {
		return 0, errors.New("data exceeds page bounds")
	}

	n := copy(p.bytes[offset:], data)
	return n, nil
}

// Read copies data from page at the specified offset and writes it to the dst slice
func (p *Page) Read(offset int, dst []byte) int {
	return copy(dst, p.bytes[offset:])
}

// WHY ???
func (p *Page) Bytes() []byte {
	return p.bytes
}

// Size returns size of underlying page byte array
func (p *Page) Size() int {
	return len(p.bytes)
}
