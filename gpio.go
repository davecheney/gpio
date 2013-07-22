// Package gpio provides an interface to the GPIO pins on various
// embedded systems.
package gpio

type Direction string

const (
	INPUT  Direction = "in"
	OUTPUT Direction = "out"
)

// Pin represents a GPIO pin.
type Pin interface {
	// Set sets the pin to a high level.
	Set()

	// Get returns the current value of the pin as a boolean value.
	// High values are true, low values are false.
	Get() bool

	// Clear sets the pin to a low value.
	Clear()

	// Close releases the pin back to the operating system
	Close() error

	// SetDirection sets the direction of the pin.
	SetDirection(Direction)

	// Diretion returns the current direction of the pin.
	Direction() Direction

	// Err returns the last error
	Err() error
}
