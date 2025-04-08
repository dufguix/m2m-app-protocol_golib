package applayerprotocol

type Command uint8

const (
	AckCmd         Command = 0x01 // Acknowledgment
	NackCmd        Command = 0x02 // Negative acknowledgment
	KeepAliveCmd   Command = 0x03 // Keep alive signal
	StartCmd       Command = 0x04 // Start command
	StopCmd        Command = 0x05 // Stop command
	GetCmd         Command = 0x06 // Get command
	GetValueCmd    Command = 0x07 // Response of "GetCmd"
	SetCmd         Command = 0x08 // Set command
	StartStreamCmd Command = 0x09 // Start streaming command
	StopStreamCmd  Command = 0x0A // Stop streaming command
	StreamValueCmd Command = 0x0B // Stream value command
)

type State uint8

const (
	StoppedState State = iota
	StartingState
	RunningState
	KeepAliveState
)

// Exclusively used with OnStarting event
type StartingCause uint8

const (
	TriggerCause   StartingCause = iota // when Start() is called.
	KeepAliveCause                      // when keepalive state didnt received an ack over time.
	StopCmdCause                        // when running or keepalive states receive a stop cmd.
	StartCmdCause                       // when running or keepalive states receive a start cmd.
)
