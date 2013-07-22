package linux

import (
	"github.com/davecheney/gpio"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLinuxImplementsPin(t *testing.T) {

	assert.Implements(t, (*gpio.Pin)(nil), new(Pin))

}

// TODO: write some useful tests
