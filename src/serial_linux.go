//go:build linux

package main

import "golang.org/x/sys/unix"

var baudrateMap = map[int]uint32{
	0:       unix.B9600,
	50:      unix.B50,
	75:      unix.B75,
	110:     unix.B110,
	134:     unix.B134,
	150:     unix.B150,
	200:     unix.B200,
	300:     unix.B300,
	600:     unix.B600,
	1200:    unix.B1200,
	1800:    unix.B1800,
	2400:    unix.B2400,
	4800:    unix.B4800,
	9600:    unix.B9600,
	19200:   unix.B19200,
	38400:   unix.B38400,
	57600:   unix.B57600,
	115200:  unix.B115200,
	230400:  unix.B230400,
	460800:  unix.B460800,
	500000:  unix.B500000,
	576000:  unix.B576000,
	921600:  unix.B921600,
	1000000: unix.B1000000,
	1152000: unix.B1152000,
	1500000: unix.B1500000,
	2000000: unix.B2000000,
	2500000: unix.B2500000,
	3000000: unix.B3000000,
	3500000: unix.B3500000,
	4000000: unix.B4000000,
}

var databitsMap = map[int]uint32{
	0: unix.CS8,
	5: unix.CS5,
	6: unix.CS6,
	7: unix.CS7,
	8: unix.CS8,
}

const ioctlTcgetattr = unix.TCGETS
const ioctlTcsetattr = unix.TCSETS
const ioctlTcsetDrain = unix.TCSETSW
const ioctlTcflsh = unix.TCFLSH
const ioctlTcavail = unix.TIOCINQ

const tcCMSPAR = unix.CMSPAR
const tcIUCLC = unix.IUCLC
const tcCRTSCTS uint32 = unix.CRTSCTS

func (port *portDto) Available() int {
	n, err := unix.IoctlGetInt(port.handle, ioctlTcavail)
	fatalIfError(err)
	return n
}

func (port *portDto) Drain() {
	err := unix.IoctlSetTermios(port.handle, ioctlTcsetDrain, port.settings)
	fatalIfError(err)
}

func (port *portDto) Discard() {
	err := unix.IoctlSetInt(port.handle, ioctlTcflsh, unix.TCIOFLUSH)
	fatalIfError(err)
}
