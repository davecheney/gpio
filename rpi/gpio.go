package rpi

import (
	"log"
	"syscall"
	"unsafe"

	"github.com/davecheney/gpio"
)

var (
	gpfsel, gpset, gpclr, gplev []*uint32
)

func initGPIO(memfd int) {
	buf, err := syscall.Mmap(memfd, BCM2835_GPIO_BASE, BCM2835_BLOCK_SIZE, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		log.Fatalf("rpi: unable to mmap GPIO page: %v", err)
	}
	gpfsel = []*uint32{
		(*uint32)(unsafe.Pointer(&buf[BCM2835_GPFSEL0])),
		(*uint32)(unsafe.Pointer(&buf[BCM2835_GPFSEL1])),
		(*uint32)(unsafe.Pointer(&buf[BCM2835_GPFSEL2])),
		(*uint32)(unsafe.Pointer(&buf[BCM2835_GPFSEL3])),
		(*uint32)(unsafe.Pointer(&buf[BCM2835_GPFSEL4])),
		(*uint32)(unsafe.Pointer(&buf[BCM2835_GPFSEL5])),
	}
	gpset = []*uint32{
		(*uint32)(unsafe.Pointer(&buf[BCM2835_GPSET0])),
		(*uint32)(unsafe.Pointer(&buf[BCM2835_GPSET1])),
	}
	gpclr = []*uint32{
		(*uint32)(unsafe.Pointer(&buf[BCM2835_GPCLR0])),
		(*uint32)(unsafe.Pointer(&buf[BCM2835_GPCLR1])),
	}
	gplev = []*uint32{
		(*uint32)(unsafe.Pointer(&buf[BCM2835_GPLEV0])),
		(*uint32)(unsafe.Pointer(&buf[BCM2835_GPLEV1])),
	}
}

// pin represents a specalised RPi GPIO pin with fast paths for
// several operations.
type pin struct {
	gpio.Pin       // the underlying Pin implementation
	pin      uint8 // the actual pin number
}

// OpenPin returns a gpio.Pin implementation specalised for the RPi.
func OpenPin(number int, mode gpio.Mode) (gpio.Pin, error) {
	initOnce.Do(initRPi)
	p, err := gpio.OpenPin(number, mode)
	return &pin{Pin: p, pin: uint8(number)}, err
}

func (p *pin) Set() {
	offset := p.pin / 32
	shift := p.pin % 32
	*gpset[offset] = (1 << shift)
}

func (p *pin) Clear() {
	offset := p.pin / 32
	shift := p.pin % 32
	*gpclr[offset] = (1 << shift)
}

func (p *pin) Get() bool {
	offset := p.pin / 32
	shift := p.pin % 32
	return *gplev[offset]&(1<<shift) == (1 << shift)
}

func GPIOFSel(pin, mode uint8) {
	offset := pin / 10
	shift := (pin % 10) * 3
	value := *gpfsel[offset]
	mask := BCM2835_GPIO_FSEL_MASK << shift
	value &= ^uint32(mask)
	value |= uint32(mode) << shift
	*gpfsel[offset] = value & mask
}
