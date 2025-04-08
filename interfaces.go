package applayerprotocol

type ItemId uint8

const (
	// User can't use id 0 or id 255, there are reserved.
	// Inherited from C lib.
	DeletedReservedId   ItemId = 0
	AvailableReservedId ItemId = 255
)

type LocalItemI interface {
	GetId() ItemId
	OnSetCmd(bytes []uint8)
	OnGetCmd(bytes []uint8) []uint8 //pass a buffer, returns updated slice header. Goal is to avoid heap alloc.
	HasChanged() bool               // Checked by task for stream data
}

type RemoteItemI interface {
	GetId() ItemId
	OnGetValueCmd(bytes []uint8)
}

type ComSenderI interface {
	Send(bytes []uint8)
}

// callbacks for debug, logs, ...
type EventListnerI interface {
	OnStarting(cause StartingCause) // when entering in StartingState
	OnStartDone()                   // when moving from starting state to running state
	OnAck(cmd Command)
	OnNack(cmd Command)
}
