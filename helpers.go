package gpio

import (
	"os"
	"strings"
)

// MergeStrings merges many strings together.
// This exists because it is much faster than sprintf.
func MergeStrings(stringArray ...string) string {
	return strings.Join(stringArray, "")
}

// MustOpenForRead attempts to open a file handle and panics
// if an error occurs
func OpenForRead(name string) (*os.File, error) {
	return os.Open(name)
}

// MustOpenForWrite attempts to open a file handle and panics
// if an error occurs
func OpenForWrite(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_WRONLY, 0200)
}

// MustOpenForReadWrite attempts to open a file handle and panics
// if an error occurs
func OpenForReadWrite(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_RDWR, 0600)
}

// MustClose attempts to close a file handle and panics
// if an error occurs
func Close(file *os.File) error {
	return file.Close()
}
