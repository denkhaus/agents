package helper

import (
	"bytes"
	"io"
)

// ByteReader creates an io.Reader from a byte slice.
func NewByteReader(data []byte) io.Reader {
	return bytes.NewReader(data)
}

// StringPtr returns a pointer to a string.
func StringPtr(s string) *string {
	return &s
}

// BoolPtr returns a pointer to a boolean.
func BoolPtr(b bool) *bool {
	return &b
}

// IntPtr returns a pointer to an integer.
func IntPtr(i int) *int {
	return &i
}
