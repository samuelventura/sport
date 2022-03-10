# Sport

Serial elixir port with a very focused async API.

## Usage

```elixir
port = Sport.open "/dev/ttyUSB0", 9600, "8N1"
true = Sport.write port, "hello\n"
true = Sport.drain port
Sport.read port, 6, 400
true = Sport.packetc port, 0x0A
true = Sport.write port, "hello\n"
"hello\n" = Sport.receive port, 400
true = Sport.close port
flush()
```

## Goals and Scope

- Very focused async API
- Minimal polling
- Erlang safe
- Posix only

##  Read use cases

- Wait for a modbus packet to arrive as slave -> Read until silent period
- Wait for a modbus response to arrive as master -> Read until n bytes
- Wait for a nl/cr delimited instrument response -> Read until delimiter
- Wait for a nl/cr delimited barcode scan -> Read until delimiter
- Wait for an undelimited barcode scan -> Read until silent period
- Wait next n bytes from a long data transfer -> Read until n bytes

## Write use cases

- Write more data to a long transfer -> Write n bytes
- Write a modbus request as master -> Write after discarding all buffers
- Write a modbus response as slave -> Write after discarding all buffers
- Write an instrument request -> Write after discarding all buffers

## Timing

VMIN (and in general PC timers) is totally unappropiated and unreliable to implement silence based packeting. Packet length information (in a header maybe) or a delimiting char is needed to avoid polling waste.

## Why golang port and not a C NIF like [sniff](https://github.com/samuelventura/sniff)?

- NIFs are more "dangerous" and ports come with integrated async 
- C is more "dangerous" and its speed is not really required
- Needed abstractions are easier to implement in golang
- Cross compiles to any mayor and most embedded platforms
- Port requires simplified error handling
- Golang is lean and fast enough

## References

- https://www.cmrr.umn.edu/~strupp/serial.html
- http://unixwiz.net/techtips/termios-vmin-vtime.html
- https://github.com/pkg/term/blob/master/termios/termios_linux.go
- https://github.com/albenik/go-serial
- https://github.com/bugst/go-serial

## Research

- Impact of GOMAXPROCS
- Impact of runtime.GC()
- Darwin support without cgo
- Detect usb adapter removal
