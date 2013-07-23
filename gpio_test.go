package gpio_test

import (
	"os"
	"testing"
)

// test helpers

// checkRoot checks if the user is root, and skips the test if not.
func checkRoot(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("test requires root privs")
	}
}

// checkNotRoot checks that the user is NOT root, and skips the test otherwise.
func checkNotRoot(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("skipping test for root user")
	}
}
