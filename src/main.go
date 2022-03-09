package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	path := getArg(0)
	speed := getArg(1)
	config := getArg(2)
	bauds, err := strconv.Atoi(speed)
	fatalIfError(err)
	mode := &Mode{BaudRate: bauds}
	configMode(config, mode)
	port := Open(path, mode)
	scanner := bufio.NewScanner(os.Stdin)
	var resp strings.Builder
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		for i := 0; i < len(line); i++ {
			c := line[i]
			switch c {
			case 'f':
				port.Drain()
			case 'd':
				port.Discard()
			case 'w':
				nn, err := strconv.ParseUint(line[i+1:i+3], 16, 8)
				fatalIfError(err)
				data := make([]byte, nn)
				for j := range data {
					k := i + 3 + 2*j
					d, err := strconv.ParseUint(line[k:k+2], 16, 8)
					fatalIfError(err)
					data[j] = byte(d)
				}
				n := port.Write(data)
				if n < int(nn) {
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
				i += 2 + 2*int(nn)
			case 'r':
				nn, err := strconv.ParseUint(line[i+1:i+3], 16, 8)
				fatalIfError(err)
				tt, err := strconv.ParseUint(line[i+3:i+5], 16, 8)
				fatalIfError(err)
				port.Packet(uint8(nn), uint8(tt))
				data := make([]byte, nn)
				n := port.Read(data)
				resp.Reset()
				resp.WriteRune('r')
				fmt.Fprintf(&resp, "%02x", n)
				for _, d := range data[:n] {
					fmt.Fprintf(&resp, "%02x", d)
				}
				fmt.Println(resp.String())
				i += 4
			}
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
