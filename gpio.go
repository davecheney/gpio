package gpio

type Mode string

const (
	ModeInput  Mode = "in"
	ModeOutput Mode = "out"
	ModePWM         = "pwm"
)

type IRQEvent func(number int, state bool)

// Pin represents a GPIO pin.
type Pin interface {
	Mode() Mode     // gets the current pin mode
	SetMode(Mode)   // set the current pin mode
	Set()           // sets the pin state high
	Clear()         // sets the pin state low
	Close() error   // if applicable, closes the pin
	Get() bool      // returns the current pin state
	Watch(IRQEvent) // calls the function argument when an edge trigger event occurs
	Wait(bool)      // wait for pin state to match boolean argument

	Err() error // returns the last error state
}
