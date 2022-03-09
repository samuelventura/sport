//go:build darwin

package main

// #include <termios.h>
// #include <sys/ioctl.h>
// static int posix_tcflush(int fd) { return tcflush(fd, TCIOFLUSH); }
// static int posix_tcdrain(int fd) { return tcdrain(fd); }
// static int posix_available(int fd) {
//	int count = 0;
//  int r = ioctl(fd, FIONREAD, &count);
//  return r < 0 ? r : count;
// }
import "C"

import (
	"fmt"

	"golang.org/x/sys/unix"
)

var baudrateMap = map[int]uint64{
	0:      unix.B9600,
	50:     unix.B50,
	75:     unix.B75,
	110:    unix.B110,
	134:    unix.B134,
	150:    unix.B150,
	200:    unix.B200,
	300:    unix.B300,
	600:    unix.B600,
	1200:   unix.B1200,
	1800:   unix.B1800,
	2400:   unix.B2400,
	4800:   unix.B4800,
	9600:   unix.B9600,
	19200:  unix.B19200,
	38400:  unix.B38400,
	57600:  unix.B57600,
	115200: unix.B115200,
	230400: unix.B230400,
}

var databitsMap = map[int]uint64{
	0: unix.CS8,
	5: unix.CS5,
	6: unix.CS6,
	7: unix.CS7,
	8: unix.CS8,
}

const ioctlTcgetattr = unix.TIOCGETA
const ioctlTcsetattr = unix.TIOCSETA

const tcCMSPAR uint64 = 0 // may be CMSPAR or PAREXT
const tcIUCLC uint64 = 0
const tcCCTS_OFLOW uint64 = 0x00010000
const tcCRTS_IFLOW uint64 = 0x00020000
const tcCRTSCTS uint64 = (tcCCTS_OFLOW | tcCRTS_IFLOW)

func (port *portDto) Available() int {
	r := C.posix_available(C.int(port.handle))
	if r < 0 {
		err := fmt.Errorf("available failed %d", r)
		fatalIfError(err)
	}
	return int(r)
}

func (port *portDto) Drain() {
	r := C.posix_tcdrain(C.int(port.handle))
	if r < 0 {
		err := fmt.Errorf("tcdrain failed %d", r)
		fatalIfError(err)
	}
}

func (port *portDto) Discard() {
	r := C.posix_tcflush(C.int(port.handle))
	if r < 0 {
		err := fmt.Errorf("tcflush failed %d", r)
		fatalIfError(err)
	}
}
