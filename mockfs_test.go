package gpio_test

// mockfs provides a /sysfs replacement for testing

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

type mockfs struct {
	t    *testing.T
	base string
}

func chk(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

// mkfifo creates a fifo(7) to emulate the control files in sysfs
func mkfifo(t *testing.T, path string) {
	chk(t, syscall.Mkfifo(path, 0777))
}

func newmockfs(t *testing.T) *mockfs {
	base, err := ioutil.TempDir("", "mockfs")
	chk(t, err)
	export := filepath.Join(base, "export")
	mkfifo(t, export)
	unexport := filepath.Join(base, "unexport")
	mkfifo(t, unexport)

	mock := mockfs{
		t:    t,
		base: base,
	}
	return &mock
}

// destroy removes the mockfs.
func (m *mockfs) destroy() {
	chk(m.t, os.RemoveAll(m.base))
}
