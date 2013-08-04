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
	pin, err := gpio.OpenPin(gpio.GPIO22, gpio.ModeInput)
	if err != nil {
		fmt.Printf("Error opening pin! %s\n", err)
		return
	}
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
			pin.Close()
			power.Close()
			os.Exit(0)
		}
	}()

	defer pin.Close()
	err = pin.BeginWatch(gpio.EdgeFalling, func(state bool) {
		fmt.Printf("Callback: Pin %d is now %v\n\n", gpio.GPIO22, state)
	})
	if err != nil {
		fmt.Printf("Unable to watch pin: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("Now watching pin 22 on a rising edge.")

	for {
		fmt.Println("Power high")
		power.Set()
		time.Sleep(1000 * time.Millisecond)
		fmt.Println("Power low")
		power.Clear()
		time.Sleep(1000 * time.Millisecond)
	}

}
