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
	queue := make(chan []byte)
	go func() {
		reader := bufio.NewReader(os.Stdin)
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
			vmin := packet[1]
			vtime := packet[2] //tenths of a second
			port.Packet(uint8(vmin), uint8(vtime))
			data := make([]byte, vmin+1)
			data[0] = 'r'
			rdata := data[1:]
			//log.Println("vmin", vmin, "vtime", vtime)
			read := 0 //single read gets ~64 bytes
			for read < len(rdata) {
				read += port.Read(rdata[read:])
			}
			send(writer, data)
		}
	}
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

func setupLog() {
	os.Setenv("GOTRACEBACK", "all")
	w := &logWriter{}
	w.pid = os.Getpid()
	log.SetFlags(0)
	log.SetOutput(w)
}
