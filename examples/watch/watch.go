package main

import (
	"fmt"
	"github.com/davecheney/gpio"
	"os"
	"os/signal"
	"time"
)

func main() {
	// set GPIO22 to input mode
	pin, err := gpio.OpenPin(gpio.GPIO22, gpio.ModeInput)
	if err != nil {
		fmt.Printf("Error opening pin! %s\n", err)
		return
	}
	// set GPIO17 to output mode
	power, err := gpio.OpenPin(gpio.GPIO17, gpio.ModeOutput)
	if err != nil {
		fmt.Printf("Error opening pin! %s\n", err)
		return
	}

	// clean up on exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			fmt.Println("Closing pin and terminating program.")
			power.Clear()
			pin.Close()
			power.Close()
			os.Exit(0)
		}
	}()

	err = pin.BeginWatch(gpio.EdgeFalling, func() {
		fmt.Printf("Callback for %d triggered!\n\n", gpio.GPIO22)
	})
	if err != nil {
		fmt.Printf("Unable to watch pin: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("Now watching pin 22 on a falling edge.")

	for {
		fmt.Println("Setting power high")
		power.Set()
		time.Sleep(2000 * time.Millisecond)
		fmt.Println("Setting power low")
		power.Clear()
		time.Sleep(2000 * time.Millisecond)
	}

}
