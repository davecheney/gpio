package linux

import (
	"bytes"
	"github.com/davecheney/gpio/common"
	"os"
	"strings"
)

// MergeStrings merges many strings together.
// This exists because it is much faster than sprintf.
func MergeStrings(stringArray ...string) string {
	return strings.Join(stringArray, "")
}

// BytesForDirection returns the byte array to use for a given Direction
func BytesForDirection(direction common.Direction) []byte {
	switch {
	case direction == common.DirectionIn:
		return BytesDirectionIn
	case direction == common.DirectionOut:
		return BytesDirectionOut
	default:
		return nil
	}
}

// DirectionForBytes returns the Direction for a given byte array
func DirectionForBytes(byteSlice []byte) common.Direction {
	switch {
	case bytes.Compare(byteSlice, BytesDirectionIn) == 0:
		return common.DirectionIn
	case bytes.Compare(byteSlice, BytesDirectionOut) == 0:
		return common.DirectionOut
	default:
		return common.DirectionNone
	}
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
