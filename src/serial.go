package main

type Port interface {
	FD() int
	Read(p []byte) int
	Write(p []byte) int
	Packet(vmin, vtime uint8)
	Available() int
	Drain()
	Discard()
	Close()
}

type Mode struct {
	BaudRate int      // platform dependant
	DataBits int      // 7 or 8
	Parity   Parity   // None, Odd and Even
	StopBits StopBits // 1, 2
}

type Parity int

const (
	NoParity Parity = iota
	OddParity
	EvenParity
)

type StopBits int

const (
	OneStopBit StopBits = iota
	TwoStopBits
)
