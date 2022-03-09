package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

func main() {
	//erlang forwards golang stderr to its own
	SetupLog()
	path := getArg(1)
	speed := getArg(2)
	config := getArg(3)
	bauds, err := strconv.Atoi(speed)
	fatalIfError(err)
	mode := &Mode{BaudRate: bauds}
	configMode(config, mode)
	port := Open(path, mode)
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	head := make([]byte, 2)
	for {
		runtime.GC()
		n, err := reader.Read(head)
		if err == nil && n != 2 {
			err = fmt.Errorf("stdin read failed %d", n)
		}
		fatalIfError(err)
		size := int(binary.BigEndian.Uint16(head))
		packet := make([]byte, size)
		n, err = reader.Read(packet)
		if err == nil && n != size {
			err = fmt.Errorf("stdin read failed %d", n)
		}
		fatalIfError(err)
		cmd := packet[0]
		switch cmd {
		case 'd':
			port.Drain()
		case 'D':
			port.Discard()
		case 'w':
			data := packet[1:]
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
		case 'r':
			vmin := packet[1]
			vtime := packet[2] //tenths of a second
			port.Packet(uint8(vmin), uint8(vtime))
			data := make([]byte, vmin)
			read := port.Read(data)
			binary.BigEndian.PutUint16(head, uint16(read))
			n, err = writer.Write(head)
			if err == nil && n != 2 {
				err = fmt.Errorf("stdout write failed %d", n)
			}
			fatalIfError(err)
			n, err = writer.Write(data[:read])
			if err == nil && n != read {
				err = fmt.Errorf("stdout write failed %d", read)
			}
			fatalIfError(err)
			err = writer.Flush()
			fatalIfError(err)
		}
	}
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

func SetupLog() {
	os.Setenv("GOTRACEBACK", "all")
	w := &logWriter{}
	w.pid = os.Getpid()
	log.SetFlags(0)
	log.SetOutput(w)
}
