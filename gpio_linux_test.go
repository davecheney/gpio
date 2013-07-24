package gpio_test

import (
	"testing"

	"github.com/davecheney/gpio"
	"github.com/davecheney/gpio/rpi"
)

func TestOpenPin(t *testing.T) {
	checkRoot(t)
	pin, err := gpio.OpenPin(rpi.GPIO_P1_22, gpio.ModeInput)
	if err != nil {
		t.Fatal(err)
	}
	err = pin.Close()
	if err != nil {
		t.Fatal(err)
	}
}

// test opening pins from a non priviledged user fails
func TestOpenPinUnpriv(t *testing.T) {
	checkNotRoot(t)
	pin, err := gpio.OpenPin(rpi.GPIO_P1_22, gpio.ModeInput)
	if err == nil {
		pin.Close()
		t.Fatalf("OpenPin is expected to fail for non priv user")
	}
}

func TestSetDirection(t *testing.T) {
	checkRoot(t)
	pin, err := gpio.OpenPin(rpi.GPIO_P1_22, gpio.ModeInput)
	if err != nil {
		t.Fatal(err)
	}
	defer pin.Close()
	if dir, err := pin.Mode(), pin.Err(); dir != gpio.ModeInput || err != nil {
		t.Fatalf("pin.Mode(): expected %v %v , got %v %v", gpio.ModeInput, nil, dir, err)
	}
	pin.SetMode(gpio.ModeOutput)
	if pin.Err() != nil {
		t.Fatal(err)
	}
	if dir, err := pin.Mode(), pin.Err(); dir != gpio.ModeOutput || err != nil {
		t.Fatalf("pin.Mode(): expected %v %v , got %v %v", gpio.ModeOutput, nil, dir, err)
	}
}
