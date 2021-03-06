package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"golang.org/x/sys/unix"
)

//erlang forwards golang stderr to its own
//which is very convenient for debuging port code
func main() {
	setupLog()
	path := getArg(1)
	speed := getArg(2)
	config := getArg(3)
	bauds, err := strconv.Atoi(speed)
	fatalIfError(err)
	mode := &Mode{BaudRate: bauds}
	configMode(config, mode)
	port := Open(path, mode)
	go func() {
		poll_close(port)
	}()
	queue := make(chan []byte)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		head := make([]byte, 2)
		for {
			runtime.GC()
			n, err := reader.Read(head)
			if err == io.EOF {
				os.Exit(1)
			}
			if err == nil && n != 2 {
				err = fmt.Errorf("stdin read failed %d", n)
			}
			fatalIfError(err)
			size := int(binary.BigEndian.Uint16(head))
			packet := make([]byte, size)
			n, err = reader.Read(packet)
			if err == io.EOF {
				os.Exit(1)
			}
			if err == nil && n != size {
				err = fmt.Errorf("stdin read failed %d", n)
			}
			fatalIfError(err)
			queue <- packet
		}
	}()
	writer := bufio.NewWriter(os.Stdout)
	for packet := range queue {
		cmd := packet[0]
		switch cmd {
		case 'd':
			sync := packet[1]
			port.Drain()
			if sync > 0 {
				send(writer, packet[:1])
			}
		case 'D':
			sync := packet[1]
			port.Discard()
			if sync > 0 {
				send(writer, packet[:1])
			}
		case 'w':
			sync := packet[1]
			data := packet[2:]
			n := port.Write(data)
			if n < len(data) {
				//write should not require to parse a response
				//unix.write is expected to block and only
				//write less then requested if there is
				//a really extreme IO error (like no HD space)
				err := fmt.Errorf("write failed %d", n)
				fatalIfError(err)
				//and sleep (tied to the baudrate) is required
				//for a proper retry since drain waits until
				//all data is sent not until enough buffer
				//is available for pending data which is
				//inefficient
			}
			if sync > 0 {
				send(writer, packet[:1])
			}
		case 'r':
			size := int(binary.BigEndian.Uint16(packet[1:3]))
			vtime := packet[3] //tenths of a second
			if size == 0 {
				size = port.Available()
			}
			read_n(port, writer, size, vtime)
		case 'p':
			prof := packet[1]
			switch prof {
			case 'n':
				size := int(binary.BigEndian.Uint16(packet[2:4]))
				go func() {
					for {
						runtime.GC()
						read_n(port, writer, size, 0)
					}
				}()
			case 'c':
				c := packet[2]
				go func() {
					for {
						runtime.GC()
						read_c(port, writer, c)
					}
				}()
			}
		}
	}
}

func poll_close(port Port) {
	fd := unix.PollFd{
		Fd:     int32(port.FD()),
		Events: unix.POLLHUP,
	}
	fds := []unix.PollFd{fd}
	_, err := unix.Poll(fds, -1)
	fatalIfError(err)
	os.Exit(2)
}

func read_c(port Port, writer *bufio.Writer, c byte) {
	var buff bytes.Buffer
	buff.WriteByte('r')
	port.Packet(1, 0)
	data := make([]byte, 1)
	for {
		n := port.Read(data)
		if n == 1 {
			cc := data[0]
			buff.WriteByte(cc)
			if c == cc {
				send(writer, buff.Bytes())
				return
			}
		}
	}
}

func read_n(port Port, writer *bufio.Writer, size int, vtime uint8) {
	data := make([]byte, size+1)
	data[0] = 'r'
	buff := data[1:]
	read := 0 //single read gets ~64 bytes
	for read < len(buff) {
		vmin := 255
		pend := len(buff) - read
		if pend < 255 {
			vmin = pend
		}
		port.Packet(uint8(vmin), vtime)
		read += port.Read(buff[read:])
	}
	send(writer, data)
}

func send(writer *bufio.Writer, packet []byte) {
	size := len(packet)
	head := make([]byte, 2)
	binary.BigEndian.PutUint16(head, uint16(size))
	n, err := writer.Write(head)
	if err == nil && n != 2 {
		err = fmt.Errorf("stdout write failed %d", n)
	}
	fatalIfError(err)
	n, err = writer.Write(packet[:size])
	if err == nil && n != size {
		err = fmt.Errorf("stdout write failed %d", size)
	}
	fatalIfError(err)
	err = writer.Flush()
	fatalIfError(err)
}

func getArg(index int) string {
	if index >= len(os.Args) {
		err := fmt.Errorf("arg not found %d", index)
		fatalIfError(err)
	}
	return os.Args[index]
}

func fatalIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func configMode(config string, mode *Mode) {
	switch config {
	case "8N1":
		mode.DataBits = 8
		mode.Parity = NoParity
		mode.StopBits = OneStopBit
	case "8E1":
		mode.DataBits = 8
		mode.Parity = EvenParity
		mode.StopBits = OneStopBit
	case "8O1":
		mode.DataBits = 8
		mode.Parity = OddParity
		mode.StopBits = OneStopBit
	case "7E1":
		mode.DataBits = 7
		mode.Parity = EvenParity
		mode.StopBits = OneStopBit
	case "7N1":
		mode.DataBits = 7
		mode.Parity = NoParity
		mode.StopBits = OneStopBit
	case "7O1":
		mode.DataBits = 7
		mode.Parity = OddParity
		mode.StopBits = OneStopBit
	case "8N2":
		mode.DataBits = 8
		mode.Parity = NoParity
		mode.StopBits = TwoStopBits
	case "8E2":
		mode.DataBits = 8
		mode.Parity = EvenParity
		mode.StopBits = TwoStopBits
	case "8O2":
		mode.DataBits = 8
		mode.Parity = OddParity
		mode.StopBits = TwoStopBits
	case "7N2":
		mode.DataBits = 7
		mode.Parity = NoParity
		mode.StopBits = TwoStopBits
	case "7E2":
		mode.DataBits = 7
		mode.Parity = EvenParity
		mode.StopBits = TwoStopBits
	case "7O2":
		mode.DataBits = 7
		mode.Parity = OddParity
		mode.StopBits = TwoStopBits
	default:
		err := fmt.Errorf("invalid config %s", config)
		fatalIfError(err)
	}
}

type logWriter struct {
	pid int
}

func (w logWriter) Write(bytes []byte) (int, error) {
	ts := time.Now().Format("20060102T150405.000")
	line := fmt.Sprintf("%s %d %s", ts, w.pid, string(bytes))
	return fmt.Fprint(os.Stderr, line)
}

func setupLog() {
	os.Setenv("GOTRACEBACK", "all")
	w := &logWriter{}
	w.pid = os.Getpid()
	log.SetFlags(0)
	log.SetOutput(w)
}
