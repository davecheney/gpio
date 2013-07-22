package linux

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
	ExportPath          = "/sys/class/gpio/export"
	UnexportPath        = "/sys/class/gpio/unexport"
	GPIOPathPrefix      = "/sys/class/gpio/gpio"
	DirectionPathSuffix = "/direction"
	ValuePathSuffix     = "/value"
)

var (
	BytesDirectionIn  = []byte{'i', 'n'}
	BytesDirectionOut = []byte{'o', 'u', 't'}
	BytesSet          = []byte{'1'}
	BytesClear        = []byte{'0'}
)
