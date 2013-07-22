package gpio

import (
	"strconv"
)

// pin represents a GPIO pin.
type pin struct {
	number        int    // the pin number
	numberAsBytes []byte // the pin number as a byte array to avoid converting each time
	modePath      string // the path to the /direction FD to avoid string joining each time
	valuePath     string // the path to the /value FD to avoid string joining each time
	err           error  //the last error
}

func OpenPin(number int) (Pin, error) {
	numString := strconv.Itoa(number)
	p := &pin{
		number:        number,
		numberAsBytes: []byte(numString),
		modePath:      MergeStrings(gpioPathPrefix, numString, modePathSuffix),
		valuePath:     MergeStrings(gpioPathPrefix, numString, valuePathSuffix),
	}

	// export this pin to create the virtual files on the system
	err := write(p.numberAsBytes, exportPath)
	p.err = err

	return p, err

}

func write(bytes []byte, path string) error {
	file, err := OpenForWrite(path)
	if err != nil {
		return err
	}

	file.Write(bytes)

	return Close(file)
}

func read(bytes *[]byte, path string) error {
	file, err := OpenForRead(path)
	if err != nil {
		return err
	}

	file.Read(*bytes)

	return Close(file)
}

// Close destroys the virtual FD on the filesystem
func (p *pin) Close() error {
	err := write(p.numberAsBytes, unexportPath)
	p.err = err
	return err
}

func (p *pin) Mode() Mode {
	bytes := make([]byte, 3)
	p.err = read(&bytes, p.modePath)
	return Mode(bytes)
}

func (p *pin) SetMode(mode Mode) {
	p.err = write([]byte(mode), p.modePath)
}

// Set will be replaced by a /dev/mem access
func (p *pin) Set() {
	p.err = write(BytesSet, p.valuePath)
}

// Clear will be replaced by a /dev/mem access
func (p *pin) Clear() {
	p.err = write(BytesClear, p.valuePath)
}

// Get will be replaced by a /dev/mem access
func (p *pin) Get() bool {
	bytes := make([]byte, 1)
	p.err = read(&bytes, p.valuePath)

	return bytes[0] != 0
}

func (p *pin) Watch(callback IRQEvent) {
	panic("Watch is not yet implemented!")
}

func (p *pin) Wait(condition bool) {
	panic("Wait is not yet implemented!")
}

func (p *pin) Err() error {
	return p.err
}
