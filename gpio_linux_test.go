package gpio_test

import (
	"testing"

	"github.com/davecheney/gpio"
	"github.com/davecheney/gpio/rpi"
)

func TestOpenPin(t *testing.T) {
	pin, err := gpio.OpenPin(rpi.GPIO_P1_22)
	if err != nil {
		t.Fatal(err)
	}
	err = pin.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetDirection(t *testing.T) {
	pin, err := gpio.OpenPin(rpi.GPIO_P1_22)
	if err != nil {
		t.Fatal(err)
	}
	defer pin.Close()
	pin.SetDirection(gpio.OUTPUT)
	if pin.Err() != nil {
		t.Fatal(err)
	}
	dir := pin.Direction()
	if err := pin.Err(); dir != gpio.OUTPUT || err != nil {
		t.Fatalf("pin.Direction(): expected %v %v , got %v %v", gpio.OUTPUT, nil, dir, err)
	}
}
