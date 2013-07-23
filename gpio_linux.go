package gpio

import (
	"fmt"
	"os"
	"strconv"
)

//Constants

//By default, pins 14 and 15 boot to UART mode, so they are going to be ignored for now.
//We can add them in later as necessary.

const (
	GPIO0  = 0
	GPIO1  = 1
	GPIO2  = 2
	GPIO3  = 3
	GPIO4  = 4
	GPIO7  = 7
	GPIO8  = 8
	GPIO9  = 9
	GPIO10 = 10
	GPIO11 = 11
	GPIO17 = 17
	GPIO18 = 18
	GPIO22 = 22
	GPIO23 = 23
	GPIO24 = 24
	GPIO25 = 25
)

const (
	exportPath      = "/sys/class/gpio/export"
	unexportPath    = "/sys/class/gpio/unexport"
	gpioPathPrefix  = "/sys/class/gpio/gpio"
	modePathSuffix  = "/direction"
	valuePathSuffix = "/value"
)

var (
	BytesSet   = []byte{'1'}
	BytesClear = []byte{'0'}
)

// pin represents a GPIO pin.
type pin struct {
	number        int      // the pin number
	numberAsBytes []byte   // the pin number as a byte array to avoid converting each time
	modePath      string   // the path to the /direction FD to avoid string joining each time
	valueFile     *os.File // the file handle for the value file
	err           error    //the last error
}

// OpenPin exports the pin, creating the virtual files necessary for interacting with the pin.
// It also sets the mode for the pin, making it ready for use.
func OpenPin(number int, mode Mode) (Pin, error) {
	numString := strconv.Itoa(number)
	p := &pin{
		number:        number,
		numberAsBytes: []byte(numString),
		modePath:      fmt.Sprintf("%s%s%s", gpioPathPrefix, numString, modePathSuffix),
	}

	// export this pin to create the virtual files on the system
	exportFile, err := os.OpenFile(exportPath, os.O_WRONLY, 0200)
	if err != nil {
		return nil, err
	}
	defer exportFile.Close()
	if _, err := exportFile.Write(p.numberAsBytes); err != nil {
		return nil, err
	}
	p.SetMode(mode)
	valueFile, err := os.Create(fmt.Sprintf("%s%s%s", gpioPathPrefix, numString, valuePathSuffix))
	p.valueFile = valueFile
	p.err = err

	return p, nil

}

// write opens a file for writing, writes the byte slice to it and closes the
// file.
func write(bytes []byte, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	file.Write(bytes)

	return file.Close()
}

// read opens a file for reading, reads the bytes slice from it and closes the file.
func read(bytes *[]byte, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	file.Read(*bytes)

	return file.Close()
}

// Close destroys the virtual files on the filesystem, unexporting the pin.
func (p *pin) Close() error {
	p.valueFile.Close()
	err := write(p.numberAsBytes, unexportPath)
	p.err = err
	return err
}

// Mode retrieves the current mode of the pin.
func (p *pin) Mode() Mode {
	bytes := make([]byte, 3)
	p.err = read(&bytes, p.modePath)
	return Mode(bytes)
}

// SetMode sets the mode of the pin.
func (p *pin) SetMode(mode Mode) {
	p.err = write([]byte(mode), p.modePath)
}

// Set sets the pin level high.
func (p *pin) Set() {
	p.valueFile.Write(BytesSet)
}

// Clear sets the pin level low.
func (p *pin) Clear() {
	p.valueFile.Write(BytesClear)
}

// Get retrieves the current pin level.
func (p *pin) Get() bool {
	bytes := make([]byte, 1)
	p.valueFile.Read(bytes)

	return bytes[0] != 0
}

// Watch waits for the edge level to be triggered and then calls the callback
func (p *pin) Watch(callback IRQEvent) {
	panic("Watch is not yet implemented!")
}

// Wait blocks while waits for the pin state to match the condition, then returns.
func (p *pin) Wait(condition bool) {
	panic("Wait is not yet implemented!")
}

// Err returns the last error encountered.
func (p *pin) Err() error {
	return p.err
}
