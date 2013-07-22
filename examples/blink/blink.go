package main

import (
	"fmt"
	"github.com/davecheney/gpio"
	"os"
	"os/signal"
	"time"
)

func main() {
	// set GPIO25 to output mode
	pin := gpio.OpenPin(25)
	pin.SetMode(gpio.ModeOutput)

	// turn the led off on exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			fmt.Printf("\nClearing and unexporting the pin.\n")
			pin.Clear()
			pin.Close()
			os.Exit(0)
		}
	}()

	for {
		pin.Set()
		time.Sleep(100 * time.Millisecond)
		pin.Clear()
		time.Sleep(100 * time.Millisecond)
	}
}
