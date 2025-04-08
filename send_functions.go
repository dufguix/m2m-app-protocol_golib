package applayerprotocol

// This file contains AppProtocol struct functions.
// ONLY those related to sending bytes.

func (ap *AppProtocol) sendBytes(bytes []uint8) {
	// if ap.sender == nil {
	// 	logMsg("Can't send. nil sender")
	// 	return
	// }
	ap.sender.Send(bytes)
	ap.sentCmdCounter++
}

func (ap *AppProtocol) sendCommand(cmd Command) {
	bytes := []uint8{uint8(cmd)}
	ap.sendBytes(bytes)
}

func (ap *AppProtocol) sendCommandWithArgs(cmd Command, bytes []uint8) {
	fullLength := len(bytes) + 1
	fullBytes := make([]uint8, fullLength)
	fullBytes[0] = uint8(cmd)
	copy(fullBytes[1:], bytes)
	ap.sendBytes(fullBytes)
}

func (ap *AppProtocol) sendCommandWithItemAndArgs(cmd Command, id ItemId, bytes []uint8) {
	if ap.state == StartingState {
		return
	}
	if id == DeletedReservedId || id == AvailableReservedId {
		return
	}
	fullLength := len(bytes) + 2
	fullBytes := make([]uint8, fullLength)
	fullBytes[0] = uint8(cmd)
	fullBytes[1] = uint8(id)
	copy(fullBytes[2:], bytes)
	ap.sendBytes(fullBytes)
}

func (ap *AppProtocol) sendAck(concernedCmd Command) {
	ap.sendCommandWithArgs(AckCmd, []uint8{uint8(concernedCmd)})
	ap.sentAckCounter++
}

func (ap *AppProtocol) sendNack(concernedCmd Command) {
	ap.sendCommandWithArgs(NackCmd, []uint8{uint8(concernedCmd)})
	ap.sentNackCounter++
}

func (ap *AppProtocol) sendKeepAlive() {
	ap.sendCommand(KeepAliveCmd)
}

func (ap *AppProtocol) sendStart() {
	ap.sendCommandWithArgs(StartCmd, []uint8{uint8(ap.version >> 8), uint8(ap.version)})
}

func (ap *AppProtocol) sendStop() {
	ap.sendCommand(StopCmd)
}

func (ap *AppProtocol) SendGet(id ItemId) {
	ap.sendCommandWithItemAndArgs(GetCmd, id, nil)
}

func (ap *AppProtocol) sendGetValue(id ItemId, valueBytes []uint8) {
	ap.sendCommandWithItemAndArgs(GetValueCmd, id, valueBytes)
}

func (ap *AppProtocol) SendSet(id ItemId, bytes []uint8) {
	ap.sendCommandWithItemAndArgs(SetCmd, id, bytes)
}

func (ap *AppProtocol) SendStartStream(id ItemId) {
	ap.sendCommandWithItemAndArgs(StartStreamCmd, id, nil)
}

func (ap *AppProtocol) SendStopStream(id ItemId) {
	ap.sendCommandWithItemAndArgs(StopStreamCmd, id, nil)
}

func (ap *AppProtocol) sendStreamValue(id ItemId, valueBytes []uint8) {
	ap.sendCommandWithItemAndArgs(StreamValueCmd, id, valueBytes)
}
