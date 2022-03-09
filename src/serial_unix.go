//go:build linux || darwin

package main

import (
	"fmt"

	"golang.org/x/sys/unix"
)

type portDto struct {
	settings *unix.Termios
	name     string
	handle   int
}

func Open(portName string, mode *Mode) (port *portDto) {
	h, err := unix.Open(portName,
		unix.O_RDWR|unix.O_NOCTTY|unix.O_NDELAY,
		0)
	fatalIfError(err)

	port = &portDto{
		handle: h,
		name:   portName,
	}

	err = unix.SetNonblock(port.handle, false)
	fatalIfError(err)

	port.settings, err = unix.IoctlGetTermios(port.handle, ioctlTcgetattr)
	fatalIfError(err)

	// parity
	switch mode.Parity {
	case NoParity:
		port.settings.Cflag &^= unix.PARENB
		port.settings.Cflag &^= unix.PARODD
		port.settings.Cflag &^= tcCMSPAR
		port.settings.Iflag &^= unix.INPCK
	case OddParity:
		port.settings.Cflag |= unix.PARENB
		port.settings.Cflag |= unix.PARODD
		port.settings.Cflag &^= tcCMSPAR
		port.settings.Iflag |= unix.INPCK
	case EvenParity:
		port.settings.Cflag |= unix.PARENB
		port.settings.Cflag &^= unix.PARODD
		port.settings.Cflag &^= tcCMSPAR
		port.settings.Iflag |= unix.INPCK
	default:
		err := fmt.Errorf("invalid parity")
		fatalIfError(err)
	}

	// baudrate
	baudrate, ok := baudrateMap[mode.BaudRate]
	if !ok {
		err := fmt.Errorf("invalid speed %d", mode.BaudRate)
		fatalIfError(err)
	}
	for _, rate := range baudrateMap {
		port.settings.Cflag &^= rate
	}
	port.settings.Cflag |= baudrate
	port.settings.Ispeed = baudrate
	port.settings.Ospeed = baudrate

	// databits
	databits, ok := databitsMap[mode.DataBits]
	if !ok {
		err := fmt.Errorf("invalid databits %d", mode.DataBits)
		fatalIfError(err)
	}
	port.settings.Cflag &^= unix.CSIZE
	port.settings.Cflag |= databits

	// stopbits
	switch mode.StopBits {
	case OneStopBit:
		port.settings.Cflag &^= unix.CSTOPB
	case TwoStopBits:
		port.settings.Cflag |= unix.CSTOPB
	default:
		err := fmt.Errorf("invalid stopbits %d", mode.StopBits)
		fatalIfError(err)
	}

	// raw mode
	port.settings.Cflag &^= tcCRTSCTS
	port.settings.Cflag |= unix.CREAD
	port.settings.Cflag |= unix.CLOCAL

	port.settings.Lflag &^= unix.ICANON
	port.settings.Lflag &^= unix.ECHO
	port.settings.Lflag &^= unix.ECHOE
	port.settings.Lflag &^= unix.ECHOK
	port.settings.Lflag &^= unix.ECHONL
	port.settings.Lflag &^= unix.ECHOCTL
	port.settings.Lflag &^= unix.ECHOPRT
	port.settings.Lflag &^= unix.ECHOKE
	port.settings.Lflag &^= unix.ISIG
	port.settings.Lflag &^= unix.IEXTEN

	port.settings.Iflag &^= unix.IXON
	port.settings.Iflag &^= unix.IXOFF
	port.settings.Iflag &^= unix.IXANY
	port.settings.Iflag &^= unix.INPCK
	port.settings.Iflag &^= unix.IGNPAR
	port.settings.Iflag &^= unix.PARMRK
	port.settings.Iflag &^= unix.ISTRIP
	port.settings.Iflag &^= unix.IGNBRK
	port.settings.Iflag &^= unix.BRKINT
	port.settings.Iflag &^= unix.INLCR
	port.settings.Iflag &^= unix.IGNCR
	port.settings.Iflag &^= unix.ICRNL
	port.settings.Iflag &^= tcIUCLC

	port.settings.Oflag &^= unix.OPOST
	port.settings.Oflag &^= unix.ONLCR
	port.settings.Oflag &^= unix.OCRNL

	port.settings.Cc[unix.VMIN] = 1
	port.settings.Cc[unix.VTIME] = 0

	err = unix.IoctlSetTermios(port.handle, ioctlTcsetattr, port.settings)
	fatalIfError(err)
	return
}

func (port *portDto) Packet(vmin, vtime uint8) {
	port.settings.Cc[unix.VMIN] = vmin
	port.settings.Cc[unix.VTIME] = vtime
	err := unix.IoctlSetTermios(port.handle, ioctlTcsetattr, port.settings)
	fatalIfError(err)
}

func (port *portDto) Close() {
	err := unix.Close(port.handle)
	fatalIfError(err)
}

func (port *portDto) Read(p []byte) (n int) {
	n, err := unix.Read(port.handle, p)
	fatalIfError(err)
	if n < 0 {
		err = fmt.Errorf("invalid read %d", n)
		fatalIfError(err)
	}
	return
}

func (port *portDto) Write(p []byte) (n int) {
	n, err := unix.Write(port.handle, p)
	fatalIfError(err)
	if n < 0 {
		err = fmt.Errorf("invalid write %d", n)
		fatalIfError(err)
	}
	return
}
