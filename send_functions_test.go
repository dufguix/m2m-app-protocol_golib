package applayerprotocol

import (
	"encoding/hex"
	"testing"
)

func TestSendAck(t *testing.T) {
	var wantedOutput = []uint8{0x01, 0x08}
	var resultOutput []uint8
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			resultOutput = u
		},
	}
	ap := AppProtocol{sender: &fakeSender}

	// Test begin here
	ap.sendAck(SetCmd)

	if !testEq(resultOutput, wantedOutput) {
		println(hex.EncodeToString(wantedOutput))
		println(hex.EncodeToString(resultOutput))
		t.Errorf("arrays dont match")
	}

	if ap.sentAckCounter != 1 {
		t.Errorf("counter was not incremented")
	}
}

func TestSendNack(t *testing.T) {
	var wantedResult = []uint8{0x02, 0x08}
	var outputResult []uint8
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			outputResult = u
		},
	}
	ap := AppProtocol{sender: &fakeSender}

	// Test begin here
	ap.sendNack(SetCmd)

	if !testEq(outputResult, wantedResult) {
		t.Errorf("arrays dont match")
	}
	if ap.sentNackCounter != 1 {
		t.Errorf("counter was not incremented")
	}
}

func TestSendKeepAlive(t *testing.T) {
	var wantedResult = []uint8{0x03}
	var outputResult []uint8
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			outputResult = u
		},
	}
	ap := AppProtocol{sender: &fakeSender}

	// Test begin here
	ap.sendKeepAlive()
	if !testEq(outputResult, wantedResult) {
		t.Errorf("arrays dont match")
	}
}

func TestSendStart(t *testing.T) {
	var wantedResult = []uint8{0x04, 0x0A, 0x01}
	var outputResult []uint8
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			outputResult = u
		},
	}
	ap := AppProtocol{version: 0x0A01, sender: &fakeSender}

	// Test begin here
	ap.sendStart()
	if !testEq(outputResult, wantedResult) {
		t.Errorf("arrays dont match")
	}
}

func TestSendStop(t *testing.T) {
	var wantedResult = []uint8{0x05}
	var outputResult []uint8
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			outputResult = u
		},
	}
	ap := AppProtocol{sender: &fakeSender}

	// Test begin here
	ap.sendStop()
	if !testEq(outputResult, wantedResult) {
		t.Errorf("arrays dont match")
	}
}

func TestSendGet(t *testing.T) {
	var wantedResult = []uint8{0x06, 0x01}
	var outputResult []uint8
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			outputResult = u
		},
	}
	ap := AppProtocol{sender: &fakeSender, state: RunningState}

	// Test begin here
	ap.SendGet(ItemId(1))
	if !testEq(outputResult, wantedResult) {
		t.Errorf("arrays dont match")
	}
}

func TestSendGetValue(t *testing.T) {
	var wantedResult = []uint8{0x07, 0x01, 0x0A, 0x0B, 0x0C}
	var outputResult []uint8
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			outputResult = u
		},
	}
	ap := AppProtocol{sender: &fakeSender, state: RunningState}

	// Test begin here
	ap.sendGetValue(ItemId(1), []uint8{0x0A, 0x0B, 0x0C})
	if !testEq(outputResult, wantedResult) {
		t.Errorf("arrays dont match")
	}
}

func TestSendSet(t *testing.T) {
	var wantedResult = []uint8{0x08, 0x01, 0x0A, 0x0B, 0x0C}
	var outputResult []uint8
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			outputResult = u
		},
	}
	ap := AppProtocol{sender: &fakeSender, state: RunningState}

	// Test begin here
	ap.SendSet(ItemId(1), []uint8{0x0A, 0x0B, 0x0C})
	if !testEq(outputResult, wantedResult) {
		t.Errorf("arrays dont match")
	}
}

func TestSendStartStream(t *testing.T) {
	var wantedResult = []uint8{0x09, 0x01}
	var outputResult []uint8
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			outputResult = u
		},
	}
	ap := AppProtocol{sender: &fakeSender, state: RunningState}

	// Test begin here
	ap.SendStartStream(ItemId(1))
	if !testEq(outputResult, wantedResult) {
		t.Errorf("arrays dont match")
	}
}

func TestSendStopStream(t *testing.T) {
	var wantedResult = []uint8{0x0A, 0x01}
	var outputResult []uint8
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			outputResult = u
		},
	}
	ap := AppProtocol{sender: &fakeSender, state: RunningState}

	// Test begin here
	ap.SendStopStream(ItemId(1))
	if !testEq(outputResult, wantedResult) {
		t.Errorf("arrays dont match")
	}
}

func TestSendStreamValue(t *testing.T) {
	var wantedResult = []uint8{0x0B, 0x01, 0x0A, 0x0B, 0x0C}
	var outputResult []uint8
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			outputResult = u
		},
	}
	ap := AppProtocol{sender: &fakeSender, state: RunningState}

	// Test begin here
	ap.sendStreamValue(ItemId(1), []uint8{0x0A, 0x0B, 0x0C})
	if !testEq(outputResult, wantedResult) {
		t.Errorf("arrays dont match")
	}
}

func TestSendReservedId(t *testing.T) {
	callbackCounter := 0
	fakeSender := FakeComSender{
		callback: func(u []uint8) {
			callbackCounter++
		},
	}
	ap := AppProtocol{sender: &fakeSender, state: RunningState}
	ap.sendCommandWithItemAndArgs(StreamValueCmd, ItemId(0), []uint8{0x0A})
	ap.sendCommandWithItemAndArgs(StreamValueCmd, ItemId(255), []uint8{0x0A})
	if callbackCounter != 0 {
		t.Errorf("Reserved IDs are not filtered.")
	}
}
