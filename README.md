gpio
====

GPIO for Go

testing
-------

The `linux` implementation uses the /sys/class/gpio kernel driver. The files in this directory are owned by root, so you cannot test as a mortal user. To run your tests, try this method

    go test -c && ./gpio.test -test.v


