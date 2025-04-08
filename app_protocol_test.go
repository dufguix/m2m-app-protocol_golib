package applayerprotocol

import (
	"encoding/hex"
	"testing"
	"time"
)

type FakeComSender struct {
	callback func([]uint8)
}

func (fcs *FakeComSender) Send(bytes []uint8) {
	fcs.callback(bytes)
}

type FakeEventListner struct {
	onStartingCallback  func(StartingCause)
	onStartDoneCallback func()
	onAckCallback       func(Command)
	onNackCallback      func(Command)
}

func (fel *FakeEventListner) OnStarting(cause StartingCause) {
	fel.onStartingCallback(cause)
}

func (fel *FakeEventListner) OnStartDone() {
	fel.onStartDoneCallback()
}

func (fel *FakeEventListner) OnAck(cmd Command) {
	fel.onAckCallback(cmd)
}

func (fel *FakeEventListner) OnNack(cmd Command) {
	fel.onNackCallback(cmd)
}

func testEq(a, b []uint8) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// func printSlice(a, b []uint8) {
// 	//just use hex.EncodeToString(payload)
// 	for _, v := range a {
// 		fmt.Printf("%X", v)
// 	}
// 	print("\n")
// 	for _, v := range b {
// 		fmt.Printf("%X", v)
// 	}
// 	print("\n")
// }

func TestTask(t *testing.T) {
	// TODO
}

func TestCheckAndSendStreams(t *testing.T) {
	var wantedOutput = []uint8{0x0B, 0x02, 0x0A, 0x0B, 0x0C}
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
			Id:                 ItemId(0x01),
			HasChangedCallback: func() bool { return false },
		},
		&SimpleLocalItem{
			Id: ItemId(0x02),
			OnGetCmdCallback: func(bytes []uint8) []uint8 {
				return []uint8{0x0A, 0x0B, 0x0C}
			},
			HasChangedCallback: func() bool { return true },
		},
	}

	ap := AppProtocol{state: RunningState, sender: &fakeSender, LocalItems: fakeLocalItems}
	ap.resetSubscriptions()
	ap.AddStreamSubscription(ItemId(1))
	ap.AddStreamSubscription(ItemId(2))
	// test begin here
	ap.CheckAndSendStreams()
	if !testEq(wantedOutput, resultOutput) {
		println(hex.EncodeToString(wantedOutput))
		println(hex.EncodeToString(resultOutput))
		t.Errorf("arrays dont match")
	}
	if senderCallbackCounter != 1 {
		t.Errorf("sender was not properly called. Counter: %v", senderCallbackCounter)
	}
}

func TestSomeStateChanges(t *testing.T) {
	var senderOutput []uint8
	senderCallbackCounter := 0
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			senderCallbackCounter++
			senderOutput = u
		},
	}
	var startingCause StartingCause
	onStartingCounter := 0
	onStartDoneCounter := 0
	fakeListner := FakeEventListner{
		onStartingCallback: func(c StartingCause) {
			startingCause = c
			onStartingCounter++
		},
		onStartDoneCallback: func() {
			onStartDoneCounter++
		},
		onAckCallback:  func(c Command) {},
		onNackCallback: func(c Command) {},
	}
	resetCounterFunc := func() {
		clear(senderOutput)
		senderCallbackCounter = 0
		startingCause = 255 //255 means nothing
		onStartingCounter = 0
		onStartDoneCounter = 0
	}
	ap := AppProtocol{state: RunningState}

	//Test begin here : StoppedState
	resetCounterFunc()
	ap.Init(1, &fakeListner, &fakeSender)

	if ap.state != StoppedState {
		t.Errorf("Should be StoppedState after Init")
	}

	// Moving to StartingState
	resetCounterFunc()
	ap.Start()
	if ap.state != StartingState {
		t.Errorf("Err #ZoZn1IzjY")
	}
	if onStartingCounter != 1 {
		t.Errorf("Err #Tx20TMblI")
	}
	if startingCause != TriggerCause {
		t.Errorf("Err #4iOUO36pW")
	}
	if senderCallbackCounter != 1 {
		t.Errorf("Err #l11sTNdZl")
	}
	if !testEq(senderOutput, []uint8{0x04, 0x00, 0x01}) {
		println(hex.EncodeToString(senderOutput))
		t.Errorf("Err #SHH8XtYES")
	}

	// Wait starting timeout for sending start again
	resetCounterFunc()
	time.Sleep(ap.startingTimeout + time.Millisecond)
	ap.Task()
	if ap.state != StartingState {
		t.Errorf("Err #WKKT7vtn0")
	}
	if onStartingCounter != 0 {
		t.Errorf("Err #gAYk3md2T")
	}
	if startingCause != 255 {
		t.Errorf("Err #TxQ0TQbEy")
	}
	if senderCallbackCounter != 1 {
		t.Errorf("Err #CJ5F7XTfq")
	}
	if !testEq(senderOutput, []uint8{0x04, 0x00, 0x01}) {
		println(hex.EncodeToString(senderOutput))
		t.Errorf("Err #q8kQ2Je48")
	}

	// Wait starting timeout for sending start again
	resetCounterFunc()
	time.Sleep(ap.startingTimeout + time.Millisecond)
	ap.Task()
	if ap.state != StartingState {
		t.Errorf("Err #8pYg2bDpC")
	}
	if onStartingCounter != 0 {
		t.Errorf("Err #s6T3Xftu2")
	}
	if startingCause != 255 {
		t.Errorf("Err #BZdOYrmQC")
	}
	if senderCallbackCounter != 1 {
		t.Errorf("Err #9bPJHcWBR")
	}
	if !testEq(senderOutput, []uint8{0x04, 0x00, 0x01}) {
		println(hex.EncodeToString(senderOutput))
		t.Errorf("Err #idMLJsdwU")
	}

	// Receive start ack, moving to running state
	resetCounterFunc()
	ap.Receive([]uint8{0x01, 0x04})
	if ap.state != RunningState {
		t.Errorf("Err #WKqT7Stn0")
	}
	if onStartDoneCounter != 1 {
		t.Errorf("Err #Xt6lphKJd")
	}
	if senderCallbackCounter != 0 {
		t.Errorf("Err #MF9hqZwEn")
	}

	// Wait RunningTimeout to fall into KeepAliveState
	resetCounterFunc()
	time.Sleep(ap.runningTimeout + time.Millisecond)
	ap.Task()
	if ap.state != KeepAliveState {
		t.Errorf("Err #8pwg26DpC")
	}
	if senderCallbackCounter != 1 {
		t.Errorf("Err #r5m9xF2X8")
	}
	if !testEq(senderOutput, []uint8{0x03}) {
		println(hex.EncodeToString(senderOutput))
		t.Errorf("Err #QgcEPj6DF")
	}

	// Receive a keepalive ack, come back into runningState
	resetCounterFunc()
	ap.Receive([]uint8{0x01, 0x03})
	if ap.state != RunningState {
		t.Errorf("Err #DVOIdhFUd")
	}
	if onStartDoneCounter != 0 {
		t.Errorf("Err #TxH0TUbEy")
	}
	if senderCallbackCounter != 0 {
		t.Errorf("Err #ZM5n1qzAu")
	}

	// Wait RunningTimeout to fall into KeepAliveState again
	resetCounterFunc()
	time.Sleep(ap.runningTimeout + time.Millisecond)
	ap.Task()
	if ap.state != KeepAliveState {
		t.Errorf("Err #yOC4X2irw")
	}
	if senderCallbackCounter != 1 {
		t.Errorf("Err #On7iZNauP")
	}
	if !testEq(senderOutput, []uint8{0x03}) {
		println(hex.EncodeToString(senderOutput))
		t.Errorf("Err #F73rXBCxm")
	}

	// Wait KeepAliveTimeout to fall into StartingState
	resetCounterFunc()
	time.Sleep(ap.keepAliveTimeout + time.Millisecond)
	ap.Task()
	if ap.state != StartingState {
		t.Errorf("Err #WKcT7Ttn0")
	}
	if onStartingCounter != 1 {
		t.Errorf("Err #MFXhqZwET")
	}
	if startingCause != KeepAliveCause {
		t.Errorf("Err #3WHtZAenH")
	}
	if senderCallbackCounter != 0 {
		t.Errorf("Err #DVoIdhFUE")
	}

	// Receive start cmd, moving to running state
	resetCounterFunc()
	ap.Receive([]uint8{0x04, 0x00, 0x01})
	if ap.state != RunningState {
		t.Errorf("Err #2QDK0eIfw")
	}
	if onStartDoneCounter != 1 {
		t.Errorf("Err #YDdc4xcom")
	}
	if senderCallbackCounter != 1 {
		t.Errorf("Err #CJ6F7gTfy")
	}
	if !testEq(senderOutput, []uint8{0x01, 0x04}) {
		println(hex.EncodeToString(senderOutput))
		t.Errorf("Err #hrtGgIun1")
	}

	//Receive stop cmd, moving to starting state again
	resetCounterFunc()
	ap.Receive([]uint8{0x05})
	if ap.state != StartingState {
		t.Errorf("Err #DVGIdcFUf")
	}
	if onStartingCounter != 1 {
		t.Errorf("Err #oCsM62G3a")
	}
	if startingCause != StopCmdCause {
		t.Errorf("Err #V49jlDo6H")
	}
	if senderCallbackCounter != 1 {
		t.Errorf("Err #uB7yI3C9w")
	}
	if !testEq(senderOutput, []uint8{0x01, 0x05}) {
		println(hex.EncodeToString(senderOutput))
		t.Errorf("Err #ml8fSi8bW")
	}

	// Receive start cmd, moving to running state
	resetCounterFunc()
	ap.Receive([]uint8{0x04, 0x00, 0x01})
	if ap.state != RunningState {
		t.Errorf("Err #s6Y3XDtuR")
	}
	if onStartDoneCounter != 1 {
		t.Errorf("Err #a3HX3MaWq")
	}
	if senderCallbackCounter != 1 {
		t.Errorf("Err #s6Y3XCtuQ")
	}
	if !testEq(senderOutput, []uint8{0x01, 0x04}) {
		println(hex.EncodeToString(senderOutput))
		t.Errorf("Err #fskaYCx1Z")
	}

	//Receive start cmd during RunningState, moving to StartingState again
	resetCounterFunc()
	ap.Receive([]uint8{0x04, 0x00, 0x01})
	if ap.state != StartingState {
		t.Errorf("Err #DVCIdZFUw")
	}
	if onStartingCounter != 1 {
		t.Errorf("Err #5vqb0QHGV")
	}
	if startingCause != StartCmdCause {
		t.Errorf("Err #HcZpPhKLL")
	}
	if senderCallbackCounter != 0 {
		println(hex.EncodeToString(senderOutput))
		t.Errorf("Err #2QjK0vIff: Should not reply")
	}

}
