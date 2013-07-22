package gpio

import (
	"github.com/davecheney/gpio/common"
	"github.com/davecheney/gpio/linux"
)

type Driver int

const (
	DriverLinux Driver = iota
)

// Pin represents a GPIO pin.
type Pin interface {
	Direction() common.Direction   // gets the current pin direction
	SetDirection(common.Direction) // set the current direction
	Set()                          // sets the pin state high
	Clear()                        // sets the pin state low
	Close()                        // if applicable, closes the pin
	Get() byte                     // returns the current pin state
	Watch() <-chan bool            // returns a channel that can be used to watch state changes on the pin, edge triggered
	Wait(bool)                     // wait for pin state to match boolean argument

	Err() error // returns the last error state
}

// OpenPin retrieves a Pin object for the specified driver that can be used
// to control the pin.
func OpenPin(pin int, driver Driver) Pin {

	switch driver {
	case DriverLinux:
		return linux.OpenPin(pin)
	}

	panic("Requested driver does not exist.")

}
