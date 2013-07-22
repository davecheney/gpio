package gpio

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLinuxImplementsPin(t *testing.T) {

	assert.Implements(t, (*Pin)(nil), new(pin))

}

// TODO: write some useful tests
