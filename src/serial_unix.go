//go:build linux || darwin || freebsd || openbsd

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

func Open(portName string, mode *Mode) (port *portDto, err error) {
	h, err := unix.Open(portName,
		unix.O_RDWR|unix.O_NOCTTY|unix.O_NDELAY,
		0)
	if err != nil {
		return
	}

	port = &portDto{
		handle: h,
		name:   portName,
	}

	// prevent handle leaks
	defer func() {
		if err != nil {
			port.Close()
			port.handle = 0
		}
	}()

	err = unix.SetNonblock(h, false)
	if err != nil {
		return
	}

	settings, err := getTermSettings(port)
	if err != nil {
		return
	}
	port.settings = settings

	err = setTermSettingsBaudrate(mode.BaudRate, settings)
	if err != nil {
		return
	}
	err = setTermSettingsParity(mode.Parity, settings)
	if err != nil {
		return
	}
	err = setTermSettingsDataBits(mode.DataBits, settings)
	if err != nil {
		return
	}
	err = setTermSettingsStopBits(mode.StopBits, settings)
	if err != nil {
		return
	}

	// Set raw mode
	// disable handshake
	settings.Cflag &^= tcCRTSCTS

	// Set local mode
	settings.Cflag |= unix.CREAD
	settings.Cflag |= unix.CLOCAL

	// Set raw mode
	settings.Lflag &^= unix.ICANON
	settings.Lflag &^= unix.ECHO
	settings.Lflag &^= unix.ECHOE
	settings.Lflag &^= unix.ECHOK
	settings.Lflag &^= unix.ECHONL
	settings.Lflag &^= unix.ECHOCTL
	settings.Lflag &^= unix.ECHOPRT
	settings.Lflag &^= unix.ECHOKE
	settings.Lflag &^= unix.ISIG
	settings.Lflag &^= unix.IEXTEN

	settings.Iflag &^= unix.IXON
	settings.Iflag &^= unix.IXOFF
	settings.Iflag &^= unix.IXANY
	settings.Iflag &^= unix.INPCK
	settings.Iflag &^= unix.IGNPAR
	settings.Iflag &^= unix.PARMRK
	settings.Iflag &^= unix.ISTRIP
	settings.Iflag &^= unix.IGNBRK
	settings.Iflag &^= unix.BRKINT
	settings.Iflag &^= unix.INLCR
	settings.Iflag &^= unix.IGNCR
	settings.Iflag &^= unix.ICRNL
	settings.Iflag &^= tcIUCLC

	settings.Oflag &^= unix.OPOST

	// Block reads until at least one char is available (no timeout)
	settings.Cc[unix.VMIN] = 1
	settings.Cc[unix.VTIME] = 0

	err = setTermSettings(port, settings)
	if err != nil {
		return
	}

	return
}

func (port *portDto) Discard() (err error) {
	err = unix.IoctlSetInt(port.handle, unix.TCFLSH, unix.TCIOFLUSH)
	return
}

func (port *portDto) Close() (err error) {
	err = unix.Close(port.handle)
	return
}

func (port *portDto) Read(p []byte) (n int, err error) {
	n, err = unix.Read(port.handle, p)
	return
}

func (port *portDto) Write(p []byte) (n int, err error) {
	n, err = unix.Write(port.handle, p)
	return
}

func setTermSettingsBaudrate(speed int, settings *unix.Termios) (err error) {
	baudrate, ok := baudrateMap[speed]
	if !ok {
		err = fmt.Errorf("invalid speed %d", speed)
		return
	}
	for _, rate := range baudrateMap {
		settings.Cflag &^= rate
	}
	settings.Cflag |= baudrate
	settings.Ispeed = toTermiosSpeedType(baudrate)
	settings.Ospeed = toTermiosSpeedType(baudrate)
	return nil
}

func setTermSettingsParity(parity Parity, settings *unix.Termios) (err error) {
	switch parity {
	case NoParity:
		settings.Cflag &^= unix.PARENB
		settings.Cflag &^= unix.PARODD
		settings.Cflag &^= tcCMSPAR
		settings.Iflag &^= unix.INPCK
	case OddParity:
		settings.Cflag |= unix.PARENB
		settings.Cflag |= unix.PARODD
		settings.Cflag &^= tcCMSPAR
		settings.Iflag |= unix.INPCK
	case EvenParity:
		settings.Cflag |= unix.PARENB
		settings.Cflag &^= unix.PARODD
		settings.Cflag &^= tcCMSPAR
		settings.Iflag |= unix.INPCK
	default:
		err = fmt.Errorf("invalid parity")
	}
	return
}

func setTermSettingsDataBits(bits int, settings *unix.Termios) (err error) {
	databits, ok := databitsMap[bits]
	if !ok {
		err = fmt.Errorf("invalid databits %d", bits)
		return
	}
	settings.Cflag &^= unix.CSIZE
	settings.Cflag |= databits
	return nil
}

func setTermSettingsStopBits(bits StopBits, settings *unix.Termios) (err error) {
	switch bits {
	case OneStopBit:
		settings.Cflag &^= unix.CSTOPB
	case OnePointFiveStopBits:
		err = fmt.Errorf("invalid stopbits %d", bits)
	case TwoStopBits:
		settings.Cflag |= unix.CSTOPB
	default:
		err = fmt.Errorf("invalid stopbits %d", bits)
	}
	return
}

func getTermSettings(port *portDto) (*unix.Termios, error) {
	return unix.IoctlGetTermios(port.handle, ioctlTcgetattr)
}

func setTermSettings(port *portDto, settings *unix.Termios) error {
	return unix.IoctlSetTermios(port.handle, ioctlTcsetattr, settings)
}
