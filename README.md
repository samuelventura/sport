# Sport

Serial port elixir port with a very focused async API.

## Text API

NN, DD, and TT are hexadecimal.

- `d` discard both input and output buffers
- `f` wait output buffer is empty
- `wNNDD..DD` write n bytes of data
- `rNNTT` read n bytes with timeout (0.1s granularity)
- `p{JSON}` run async packetizer profile

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
- Golang is fast enough

## References

- https://www.cmrr.umn.edu/~strupp/serial.html
- http://unixwiz.net/techtips/termios-vmin-vtime.html
- https://github.com/pkg/term/blob/master/termios/termios_linux.go

## Research

- Impact of GOMAXPROCS
