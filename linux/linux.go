package linux

import (
	"github.com/davecheney/gpio/common"
	"strconv"
)

// Pin represents a GPIO pin.
type Pin struct {
	number        int
	numberAsBytes []byte
	directionPath string
	valuePath     string

	direction common.Direction
	lastError error
}

func OpenPin(number int) *Pin {
	numString := strconv.Itoa(number)
	p := &Pin{
		number:        number,
		numberAsBytes: []byte(numString),
		directionPath: MergeStrings(GPIOPathPrefix, numString, DirectionPathSuffix),
		valuePath:     MergeStrings(GPIOPathPrefix, numString, ValuePathSuffix),
		direction:     common.DirectionNone,
	}

	// export this pin to create the virtual files on the system
	p.lastError = write(p.numberAsBytes, ExportPath)

	return p

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
func (p *Pin) Close() {
	p.lastError = write(p.numberAsBytes, UnexportPath)
}

func (p *Pin) Direction() common.Direction {
	bytes := make([]byte, 2)
	p.lastError = read(&bytes, p.directionPath)
	return DirectionForBytes(bytes)
}

func (p *Pin) SetDirection(direction common.Direction) {
	if p.direction != direction {
		p.lastError = write(BytesForDirection(direction), p.directionPath)
		p.direction = direction
	}
}

func (p *Pin) Set() {
	p.lastError = write(BytesSet, p.valuePath)
}

func (p *Pin) Clear() {
	p.lastError = write(BytesClear, p.valuePath)
}

func (p *Pin) Get() byte {
	bytes := make([]byte, 1)
	p.lastError = read(&bytes, p.valuePath)

	return bytes[0]
}

func (p *Pin) Watch() <-chan bool {
	panic("Watch is not yet implemented!")
	return make(chan bool, 0)
}

func (p *Pin) Wait(condition bool) {
	panic("Wait is not yet implemented!")
}

func (p *Pin) Err() error {
	return p.lastError
}
