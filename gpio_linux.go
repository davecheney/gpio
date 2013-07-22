package gpio

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const gpiobase = "/sys/class/gpio"

// OpenPin requests the pin be exposed by the operating system.
func OpenPin(p int) (Pin, error) {
	pinbase, err := expose(p)
	return &pin{pin: p, pinbase: pinbase}, err
}

// pin reprents a linux GPIO user space GPIO pin.
type pin struct {
	pin     int    // gpio pin number
	pinbase string // base path to pin control files
	err     error  // last error, if set
}

func (p *pin) SetDirection(dir Direction) {
	p.err = writeFile(filepath.Join(p.pinbase, "direction"), string(dir))
}

func (p *pin) Direction() Direction {
	dir, err := readFile(filepath.Join(p.pinbase, "direction"))
	p.err = err
	return Direction(dir)
}

func (p *pin) Clear()    {}
func (p *pin) Set()      {}
func (p *pin) Get() bool { return false }

func (p *pin) Err() error { return p.err }

func (p *pin) Close() error {
	return writeFile(filepath.Join(gpiobase, "unexport"), "%d", p.pin)
}

func expose(pin int) (string, error) {
	err := writeFile(filepath.Join(gpiobase, "export"), "%d", pin)
	return filepath.Join(gpiobase, fmt.Sprintf("gpio%d", pin)), err
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
