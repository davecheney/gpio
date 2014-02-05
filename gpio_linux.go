package gpio

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

//By default, pins 14 and 15 boot to UART mode, so they are going to be ignored for now.
//We can add them in later as necessary.

const (
	GPIO0     = 0
	GPIO1     = 1
	GPIO2     = 2
	GPIO3     = 3
	GPIO4     = 4
	GPIO7     = 7
	GPIO8     = 8
	GPIO9     = 9
	GPIO10    = 10
	GPIO11    = 11
	GPIO17    = 17
	GPIO18    = 18
	GPIO22    = 22
	GPIO23    = 23
	GPIO24    = 24
	GPIO25    = 25
	GPIOCount = 16 // the number of GPIO pins available
)

const (
	gpiobase     = "/sys/class/gpio"
	exportPath   = "/sys/class/gpio/export"
	unexportPath = "/sys/class/gpio/unexport"
)

var (
	bytesSet   = []byte{'1'}
	bytesClear = []byte{'0'}
)

// watchEventCallbacks is a map of pins and their callbacks when
// watching for interrupts
var watchEventCallbacks map[int]*pin

// epollFD is the FD for epoll
var epollFD int

func init() {
	setupEpoll()
	watchEventCallbacks = make(map[int]*pin)
}

// setupEpoll sets up epoll for use
func setupEpoll() {
	var err error
	epollFD, err = syscall.EpollCreate1(0)
	if err != nil {
		fmt.Println("Unable to create epoll FD: ", err.Error())
		os.Exit(1)
	}

	go func() {

		var epollEvents [GPIOCount]syscall.EpollEvent

		for {
			numEvents, err := syscall.EpollWait(epollFD, epollEvents[:], -1)
			if err != nil {
				panic(fmt.Sprintf("EpollWait error: %s", err.Error()))
			}
			for i := 0; i < numEvents; i++ {
				if eventPin, exists := watchEventCallbacks[int(epollEvents[i].Fd)]; exists {
					if eventPin.initial {
						eventPin.initial = false
					} else {
						eventPin.callback()
					}
				}
			}
		}

	}()

}

// pin represents a GPIO pin.
type pin struct {
	number        int      // the pin number
	numberAsBytes []byte   // the pin number as a byte array to avoid converting each time
	modePath      string   // the path to the /direction FD to avoid string joining each time
	edgePath      string   // the path to the /edge FD to avoid string joining each time
	valueFile     *os.File // the file handle for the value file
	callback      IRQEvent // the callback function to call when an interrupt occurs
	initial       bool     // is this the initial epoll trigger?
	err           error    //the last error
}

// OpenPin exports the pin, creating the virtual files necessary for interacting with the pin.
// It also sets the mode for the pin, making it ready for use.
func OpenPin(n int, mode Mode) (Pin, error) {
	// export this pin to create the virtual files on the system
	pinBase, err := expose(n)
	if err != nil {
		return nil, err
	}
	value, err := os.OpenFile(filepath.Join(pinBase, "value"), os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	p := &pin{
		number:    n,
		modePath:  filepath.Join(pinBase, "direction"),
		edgePath:  filepath.Join(pinBase, "edge"),
		valueFile: value,
		initial:   true,
	}
	if err := p.setMode(mode); err != nil {
		p.Close()
		return nil, err
	}
	return p, nil
}

// write opens a file for writing, writes the byte slice to it and closes the
// file.
func write(buf []byte, path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	if _, err := file.Write(buf); err != nil {
		return err
	}
	return file.Close()
}

// read opens a file for reading, reads the bytes slice from it and closes the file.
func read(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

// Close destroys the virtual files on the filesystem, unexporting the pin.
func (p *pin) Close() error {
	return writeFile(filepath.Join(gpiobase, "unexport"), "%d", p.number)
}

// Mode retrieves the current mode of the pin.
func (p *pin) Mode() Mode {
	var mode string
	mode, p.err = readFile(p.modePath)
	return Mode(mode)
}

// SetMode sets the mode of the pin.
func (p *pin) SetMode(mode Mode) {
	p.err = p.setMode(mode)
}

func (p *pin) setMode(mode Mode) error {
	return write([]byte(mode), p.modePath)
}

// Set sets the pin level high.
func (p *pin) Set() {
	_, p.err = p.valueFile.Write(bytesSet)
}

// Clear sets the pin level low.
func (p *pin) Clear() {
	_, p.err = p.valueFile.Write(bytesClear)
}

// Get retrieves the current pin level.
func (p *pin) Get() bool {
	bytes := make([]byte, 1)
	_, p.err = p.valueFile.ReadAt(bytes, 0)
	return bytes[0] == bytesSet[0]
}

// Watch waits for the edge level to be triggered and then calls the callback
// Watch sets the pin mode to input on your behalf, then establishes the interrupt on
// the edge provided

func (p *pin) BeginWatch(edge Edge, callback IRQEvent) error {
	p.SetMode(ModeInput)
	if err := write([]byte(edge), p.edgePath); err != nil {
		return err
	}

	var event syscall.EpollEvent
	event.Events = syscall.EPOLLIN | (syscall.EPOLLET & 0xffffffff) | syscall.EPOLLPRI

	fd := int(p.valueFile.Fd())

	p.callback = callback
	watchEventCallbacks[fd] = p

	if err := syscall.SetNonblock(fd, true); err != nil {
		return err
	}

	event.Fd = int32(fd)

	if err := syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, fd, &event); err != nil {
		return err
	}

	return nil

}

// EndWatch stops watching the pin
func (p *pin) EndWatch() error {

	fd := int(p.valueFile.Fd())

	if err := syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_DEL, fd, nil); err != nil {
		return err
	}

	if err := syscall.SetNonblock(fd, false); err != nil {
		return err
	}

	delete(watchEventCallbacks, fd)

	return nil

}

// Wait blocks while waits for the pin state to match the condition, then returns.
func (p *pin) Wait(condition bool) {
	panic("Wait is not yet implemented!")
}

// Err returns the last error encountered.
func (p *pin) Err() error {
	return p.err
}

func expose(pin int) (string, error) {
	pinBase := filepath.Join(gpiobase, fmt.Sprintf("gpio%d", pin))
	var err error
	if _, statErr := os.Stat(pinBase); os.IsNotExist(statErr) {
		err = writeFile(filepath.Join(gpiobase, "export"), "%d", pin)
	}
	return pinBase, err
}

func writeFile(path string, format string, args ...interface{}) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, format, args...)
	return err
}

func readFile(path string) (string, error) {
	buf, err := ioutil.ReadFile(path)
	return strings.TrimSpace(string(buf)), err
}
