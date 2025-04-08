package applayerprotocol

import (
	"errors"
	"time"

	timer "github.com/dufguix/go_timer"
)

const StreamSubsLength = 20
const MaxOngetfuncBufferSize = 10

var LogAppProtocolFunc func(msg string) = nil

func logMsg(msg string) {
	if LogAppProtocolFunc == nil {
		return
	}
	LogAppProtocolFunc(msg)
}

// Below concerns AppProtocol struct and some of its methods.
// Check receive_functions.go and send_functions.go files for all of its methods.

type AppProtocol struct {
	version             uint16
	state               State
	sender              ComSenderI
	eventListener       EventListnerI
	sentCmdCounter      uint32
	sentAckCounter      uint32
	sentNackCounter     uint32
	receivedAckCounter  uint32
	receivedNackCounter uint32
	LocalItems          []LocalItemI
	RemoteItems         []RemoteItemI
	streamSubscriptions [StreamSubsLength]ItemId //Contains the list if item ID. In the array, 0 means deleted. 255 means end of allocated subscription. Unless I find another approach for a simple list.
	streamTimer         timer.Timer
	streamTimeout       time.Duration // ms
	stateTimer          timer.Timer
	startingTimeout     time.Duration // ms
	runningTimeout      time.Duration // ms
	keepAliveTimeout    time.Duration // ms
}

func (ap *AppProtocol) Init(version uint16, eventListener EventListnerI, comSender ComSenderI) error {
	if version == 0 {
		return errors.New("version cant be 0")
	}
	if eventListener == nil || comSender == nil {
		return errors.New("cant pass nil interfaces")
	}
	ap.version = version
	ap.sender = comSender
	ap.eventListener = eventListener
	ap.streamTimeout = 500 * time.Millisecond
	ap.startingTimeout = 300 * time.Millisecond
	ap.runningTimeout = 100 * time.Millisecond
	ap.keepAliveTimeout = 100 * time.Millisecond
	ap.resetSubscriptions()
	ap.setState(StoppedState)
	return nil
}

func (ap *AppProtocol) Start() {
	ap.setState(StartingState)
	ap.eventListener.OnStarting(TriggerCause)
	ap.sendStart()
}

func (ap *AppProtocol) Stop() {
	ap.setState(StoppedState)
	ap.sendStop()
}

func (ap *AppProtocol) Task() {
	switch ap.state {
	case RunningState:
		if ap.stateTimer.RefreshAndCheck(ap.runningTimeout) {
			ap.sendKeepAlive()
			ap.setState(KeepAliveState)
			break
		}
		if ap.streamTimer.RefreshAndCheck(ap.streamTimeout) {
			ap.CheckAndSendStreams()
			ap.streamTimer.Reset()
		}
	case StartingState:
		if ap.stateTimer.RefreshAndCheck(ap.startingTimeout) {
			ap.sendStart()
			ap.stateTimer.Reset()
		}
	case KeepAliveState:
		if ap.stateTimer.RefreshAndCheck(ap.keepAliveTimeout) {
			ap.setState(StartingState)
			ap.eventListener.OnStarting(KeepAliveCause)
			break
		}
		if ap.streamTimer.RefreshAndCheck(ap.streamTimeout) {
			ap.CheckAndSendStreams()
			ap.streamTimer.Reset()
		}
		// case StoppedState:
		// 	//
	}
}

func (ap *AppProtocol) resetSubscriptions() {
	for i := range ap.streamSubscriptions {
		ap.streamSubscriptions[i] = 255
	}
}

// returns nil if not found
func (ap *AppProtocol) GetLocalItemById(id ItemId) LocalItemI {
	for _, item := range ap.LocalItems {
		if item.GetId() == id {
			return item
		}
	}
	return nil
}

// returns nil if not found
func (ap *AppProtocol) GetRemoteItemById(id ItemId) RemoteItemI {
	for _, item := range ap.RemoteItems {
		if item.GetId() == id {
			return item
		}
	}
	return nil
}

// Returns true if success
func (ap *AppProtocol) AddStreamSubscription(id ItemId) bool {
	localItem := ap.GetLocalItemById(id)
	if localItem == nil {
		return false
	}
	for index, itemId := range ap.streamSubscriptions {
		// 0 is reserved for deleted stream sub. 255 is reserved to empty cases.
		if itemId == id || itemId == DeletedReservedId || itemId == AvailableReservedId {
			ap.streamSubscriptions[index] = id
			return true
		}
	}
	return false
}

// Returns true if success
func (ap *AppProtocol) RemoveStreamSubscription(id ItemId) bool {
	for index, itemId := range ap.streamSubscriptions {
		if itemId == 255 {
			return false
		}
		if itemId == id {
			ap.streamSubscriptions[index] = 0
			return true
		}
	}
	return false
}

func (ap *AppProtocol) CheckAndSendStreams() {
	for _, itemId := range ap.streamSubscriptions {
		if itemId == 0 {
			continue
		}
		if itemId == 255 {
			break
		}

		localItem := ap.GetLocalItemById(itemId)
		if localItem == nil {
			continue
		}
		if !localItem.HasChanged() {
			continue
		}
		resultBuffer := make([]uint8, MaxOngetfuncBufferSize)
		resultBuffer = localItem.OnGetCmd(resultBuffer)
		resultLength := len(resultBuffer)
		if resultLength < 1 {
			continue
		}
		if resultLength > MaxOngetfuncBufferSize {
			logMsg("Can't send full data. MAX_ONGETFUNC_BUFFER_SIZE is too small")
			resultLength = MaxOngetfuncBufferSize
		}
		ap.sendStreamValue(itemId, resultBuffer[:resultLength])
	}
}

func (ap *AppProtocol) setState(state State) {
	ap.state = state
	ap.stateTimer.Reset()
}
