# M2M App Layer Protocol
Protocol for machine to machine:
- microcontroller to microcontroller
- microcontroller to pc
- pc to pc

Protocol for duplex communication.
Set value, Get value, Control stream of values from both sides.
You can use this protocol over different protocols: UART, USB, TCP socket, ...

Simple parsing.
Small payload:
 - 1 byte for a command.
 - 1 byte for item number
 - 1 or more bytes for item value
 
Up to you to define set/get callbacks for all yours items.
Create an array of max 254 items for local items.
Create an array of max 254 items for remote items.

Keep alive mechanism to detect broken link.

Version number at start. Up to you to define a version strategy for compatibility.

C lib available [here](https://github.com/dufguix/m2m-app-protocol_clib).

## Install
```go
go get github.com/dufguix/m2m-app-protocol_golib

import (
	ap "github.com/dufguix/m2m-app-protocol_golib"
)
```
## TODO
- protect timers to multi-process access. and sender ? Receive and Task functions
- remove ack and nack counter ?? dev can just use onAck/onNack to count/log...
- add drawio schema
- avoid to send something during StoppedState
- add example in readme