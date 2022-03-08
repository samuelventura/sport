package main

//go:generate go run golang.org/x/sys/windows/mkwinsyscall -output zsyscall_windows.go syscall_windows.go

//This a very opinionated and focused API

//Read use cases
//- Wait for a modbus packet to arrive as slave -> Read until silent period
//- Wait for a modbus response to arrive as master -> Read until n bytes
//- Wait for a nl/cr delimited instrument response -> Read until delimiter
//- Wait for a nl/cr delimited barcode scan -> Read until delimiter
//- Wait for an undelimited barcode scan -> Read until silent period
//- Wait next n bytes from a long data transfer -> Read until n bytes

//Write use cases
//- Write more data to a long transfer -> Write n bytes
//- Write a modbus request as master -> Write after discarding all buffers
//- Write a modbus response as slave -> Write after discarding all buffers
//- Write an instrument request -> Write after discarding all buffers

type Port interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Discard() error
	Close() error
}

type Mode struct {
	BaudRate int      // platform dependant
	DataBits int      // 7 or 8
	Parity   Parity   // None, Odd and Even
	StopBits StopBits // 1, 1.5, 2
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
	OnePointFiveStopBits
	TwoStopBits
)
