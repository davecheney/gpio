package gpio

// Mode represents a state of a GPIO pin
type Mode string

const (
	ModeInput  Mode = "in"
	ModeOutput Mode = "out"
	ModePWM         = "pwm"
)

// Edge represents the edge on which a pin interrupt is triggered
type Edge string

const (
	EdgeNone    Edge = "none"
	EdgeRising  Edge = "rising"
	EdgeFalling Edge = "falling"
	EdgeBoth    Edge = "both"
)

// IRQEvent defines the callback function used to inform the caller
// of an interrupt.
type IRQEvent func()

// Pin represents a GPIO pin.
type Pin interface {
	Mode() Mode                      // gets the current pin mode
	SetMode(Mode)                    // set the current pin mode
	Set()                            // sets the pin state high
	Clear()                          // sets the pin state low
	Close() error                    // if applicable, closes the pin
	Get() bool                       // returns the current pin state
	BeginWatch(Edge, IRQEvent) error // calls the function argument when an edge trigger event occurs
	EndWatch() error                 // stops watching the pin
	Wait(bool)                       // wait for pin state to match boolean argument

	Err() error // returns the last error state
}
