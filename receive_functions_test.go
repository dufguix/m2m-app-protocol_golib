package applayerprotocol

import (
	"encoding/hex"
	"testing"
)

func TestReceiveWhenStarting(t *testing.T) {
	//TODO
}

func TestReceiveAck(t *testing.T) {
	var input = []uint8{0x01, 0x08}
	senderCallbackCounter := 0
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			senderCallbackCounter++
		},
	}
	var callbackCommand Command
	fakeListner := FakeEventListner{
		onAckCallback: func(c Command) {
			callbackCommand = c
		},
	}
	ap := AppProtocol{state: KeepAliveState, sender: &fakeSender, eventListener: &fakeListner}

	// Test begin here
	ap.Receive(input)

	if ap.receivedAckCounter != 1 {
		t.Errorf("counter was not incremented")
	}
	if senderCallbackCounter != 0 {
		t.Errorf("sender should not be called")
	}
	if callbackCommand != Command(0x08) {
		t.Errorf("eventListner was not called or command was changed during the process.")
	}
	if ap.state != RunningState {
		t.Errorf("running state should be forced during the process")
	}
}

func TestReceiveNack(t *testing.T) {
	var input = []uint8{0x02, 0x08}
	senderCallbackCounter := 0
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			senderCallbackCounter++
		},
	}
	var callbackCommand Command
	fakeListner := FakeEventListner{
		onNackCallback: func(c Command) {
			callbackCommand = c
		},
	}
	ap := AppProtocol{state: KeepAliveState, sender: &fakeSender, eventListener: &fakeListner}

	// Test begin here
	ap.Receive(input)

	if ap.receivedNackCounter != 1 {
		t.Errorf("counter was not incremented")
	}
	if senderCallbackCounter != 0 {
		t.Errorf("sender should not be called")
	}
	if callbackCommand != Command(0x08) {
		t.Errorf("eventListner was not called or command was changed during the process.")
	}
	if ap.state != RunningState {
		t.Errorf("running state should be forced during the process")
	}
}

func TestReceiveKeepalive(t *testing.T) {
	var input = []uint8{0x03}
	senderCallbackCounter := 0
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			senderCallbackCounter++
		},
	}
	ap := AppProtocol{state: KeepAliveState, sender: &fakeSender}

	// Test begin here
	ap.Receive(input)

	if senderCallbackCounter != 1 {
		t.Errorf("sender should be called")
	}
}
func TestReceiveStart(t *testing.T) {
	var input = []uint8{0x04}
	senderCallbackCounter := 0
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			senderCallbackCounter++
		},
	}
	listenerCallbackCounter := 0
	fakeListner := FakeEventListner{
		onStartingCallback: func(c StartingCause) {
			if c != StartCmdCause {
				return
			}
			listenerCallbackCounter++
		},
	}
	ap := AppProtocol{state: RunningState, sender: &fakeSender, eventListener: &fakeListner}

	// Test begin here
	ap.Receive(input)
	if listenerCallbackCounter != 1 {
		t.Errorf("listner onStarting was not called")
	}
	if senderCallbackCounter != 0 {
		t.Errorf("sender should not be called")
	}
	if ap.state != StartingState {
		t.Errorf("State should change to StartingState.")
	}
}
func TestReceiveStop(t *testing.T) {
	var input = []uint8{0x05}
	senderCallbackCounter := 0
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			senderCallbackCounter++
		},
	}
	listenerCallbackCounter := 0
	fakeListner := FakeEventListner{
		onStartingCallback: func(cause StartingCause) {
			if cause != StopCmdCause {
				return
			}
			listenerCallbackCounter++
		},
	}
	ap := AppProtocol{state: RunningState, sender: &fakeSender, eventListener: &fakeListner}

	// Test begin here
	ap.Receive(input)
	if listenerCallbackCounter != 1 {
		t.Errorf("listner OnStarting(StopCmdCause) was not called")
	}
	if senderCallbackCounter != 1 {
		t.Errorf("sender should be called. Counter: %v", senderCallbackCounter)
	}
	if ap.state != StartingState {
		t.Errorf("State should be startingState")
	}
}
func TestReceiveGet(t *testing.T) {
	var input = []uint8{0x06, 0x01}
	var wantedOutput = []uint8{0x07, 0x01, 0x0A, 0x0B, 0x0C}
	var resultOutput []uint8
	senderCallbackCounter := 0
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			senderCallbackCounter++
			resultOutput = u
		},
	}
	fakeLocalItems := []LocalItemI{
		&SimpleLocalItem{
			Id: ItemId(0x01),
			OnGetCmdCallback: func(bytes []uint8) []uint8 {
				return []uint8{0x0A, 0x0B, 0x0C}
			}},
	}
	ap := AppProtocol{state: RunningState, sender: &fakeSender, LocalItems: fakeLocalItems}

	// Test begin here
	ap.Receive(input)

	if !testEq(wantedOutput, resultOutput) {
		println(hex.EncodeToString(wantedOutput))
		println(hex.EncodeToString(resultOutput))
		t.Errorf("arrays dont match")
	}
}
func TestReceiveGetValue(t *testing.T) {
	var input = []uint8{0x07, 0x01, 0x0A, 0x0B, 0x0C}
	var wantedOutput = []uint8{0x0A, 0x0B, 0x0C}
	var resultOutput []uint8

	senderCallbackCounter := 0
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			senderCallbackCounter++
		},
	}
	fakeRemoteItems := []RemoteItemI{
		&SimpleRemoteItem{
			Id: ItemId(0x01),
			OnGetValueCmdCallback: func(bytes []uint8) {
				resultOutput = bytes
			}},
	}
	ap := AppProtocol{state: RunningState, sender: &fakeSender, RemoteItems: fakeRemoteItems}

	// Test begin here
	ap.Receive(input)
	if !testEq(wantedOutput, resultOutput) {
		println(hex.EncodeToString(wantedOutput))
		println(hex.EncodeToString(resultOutput))
		t.Errorf("arrays dont match")
	}
	if senderCallbackCounter != 1 {
		t.Errorf("sender was not properly called. Counter: %v", senderCallbackCounter)
	}
}
func TestReceiveSet(t *testing.T) {
	var input = []uint8{0x08, 0x01, 0x0A, 0x0B, 0x0C}
	var wantedOutput = []uint8{0x0A, 0x0B, 0x0C}
	var resultOutput []uint8
	senderCallbackCounter := 0
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			senderCallbackCounter++
		},
	}
	fakeLocalItems := []LocalItemI{
		&SimpleLocalItem{
			Id: ItemId(0x01),
			OnSetCmdCallback: func(bytes []uint8) {
				resultOutput = bytes
			}},
	}
	ap := AppProtocol{state: RunningState, sender: &fakeSender, LocalItems: fakeLocalItems}

	// Test begin here
	ap.Receive(input)
	if !testEq(wantedOutput, resultOutput) {
		println(hex.EncodeToString(wantedOutput))
		println(hex.EncodeToString(resultOutput))
		t.Errorf("arrays dont match")
	}
	if senderCallbackCounter != 1 {
		t.Errorf("sender was not properly called. Counter: %v", senderCallbackCounter)
	}
}
func TestReceiveStartStream(t *testing.T) {
	var input = []uint8{0x09, 0x01}
	senderCallbackCounter := 0
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			senderCallbackCounter++
		},
	}
	fakeLocalItems := []LocalItemI{
		&SimpleLocalItem{
			Id: ItemId(0x01),
		},
	}
	ap := AppProtocol{state: RunningState, sender: &fakeSender, LocalItems: fakeLocalItems}
	for i := range ap.streamSubscriptions {
		ap.streamSubscriptions[i] = 255
	}
	ap.Receive(input)
	if ap.streamSubscriptions[0] != ItemId(1) {
		t.Errorf("Stream subscription not saved.")
	}
	if ap.streamSubscriptions[1] != AvailableReservedId {
		t.Errorf("The subscription list has not been instantiated correctly.")
	}
	if senderCallbackCounter != 1 {
		t.Errorf("sender was not properly called. Counter: %v", senderCallbackCounter)
	}

}
func TestReceiveStopStream(t *testing.T) {
	var input = []uint8{0x0A, 0x01}
	senderCallbackCounter := 0
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			senderCallbackCounter++
		},
	}
	fakeLocalItems := []LocalItemI{
		&SimpleLocalItem{
			Id: ItemId(0x01),
		},
	}
	ap := AppProtocol{state: RunningState, sender: &fakeSender, LocalItems: fakeLocalItems}
	ap.resetSubscriptions()
	ap.streamSubscriptions[0] = ItemId(1)
	ap.Receive(input)
	if ap.streamSubscriptions[0] != DeletedReservedId {
		t.Errorf("Stream subscription not deleted. ItemId: %v", ap.streamSubscriptions[0])
	}
	if ap.streamSubscriptions[1] != AvailableReservedId {
		t.Errorf("The subscription list has not been instantiated correctly.")
	}
	if senderCallbackCounter != 1 {
		t.Errorf("sender was not properly called. Counter: %v", senderCallbackCounter)
	}

}
func TestReceiveStreamValue(t *testing.T) {
	var input = []uint8{0x0B, 0x01, 0x0A, 0x0B, 0x0C}
	var wantedOutput = []uint8{0x0A, 0x0B, 0x0C}
	var resultOutput []uint8
	senderCallbackCounter := 0
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			senderCallbackCounter++
		},
	}
	fakeRemoteItems := []RemoteItemI{
		&SimpleRemoteItem{
			Id: ItemId(0x01),
			OnGetValueCmdCallback: func(bytes []uint8) {
				resultOutput = bytes
			}},
	}
	ap := AppProtocol{state: RunningState, sender: &fakeSender, RemoteItems: fakeRemoteItems}
	// Test begin here
	ap.Receive(input)
	if !testEq(wantedOutput, resultOutput) {
		println(hex.EncodeToString(wantedOutput))
		println(hex.EncodeToString(resultOutput))
		t.Errorf("arrays dont match")
	}
	if senderCallbackCounter != 1 {
		t.Errorf("sender was not properly called. Counter: %v", senderCallbackCounter)
	}
}
