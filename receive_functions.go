package applayerprotocol

// This file contains functions of AppProtocol struct.
// ONLY those concerning bytes reception.

func (ap *AppProtocol) receiveWhenStarting(bytes []uint8) {
	command := Command(bytes[0])
	switch command {
	case AckCmd:
		//Check if this Ack concerns a previous Start cmd.
		if len(bytes) < 2 && Command(bytes[1]) != StartCmd {
			return
		}
	case StartCmd:
		if len(bytes) < 3 {
			ap.sendNack(StartCmd)
			return
		}
		version := uint16(bytes[1])<<8 | uint16(bytes[2])
		if ap.version != version {
			ap.sendNack(StartCmd)
			return
		}
		ap.sendAck(StartCmd)
	default:
		//Silent drop other commands
		return
	}
	ap.setState(RunningState)
	ap.eventListener.OnStartDone()
}

func (ap *AppProtocol) Receive(bytes []uint8) {
	length := len(bytes)
	if length < 1 {
		return
	}
	if ap.state <= StartingState {
		if ap.state == StartingState {
			ap.receiveWhenStarting(bytes)
		}
		//silent drop when ap.state == StoppedState
		return
	}
	// Following switch concerns RUNNING_STATE and KEEP_ALIVE_STATE
	command := Command(bytes[0])
	switch command {
	case AckCmd:
		if length < 2 {
			return
		}
		ap.receivedAckCounter++
		ap.eventListener.OnAck(Command(bytes[1]))
		// reset timer to not fall into keepaliveState
		// or comeback into runningState if already into keepaliveState
		ap.setState(RunningState)
	case NackCmd:
		if length < 2 {
			return
		}
		ap.receivedNackCounter++
		ap.eventListener.OnNack(Command(bytes[1]))
		// reset timer to not fall into keepaliveState
		// or comeback into runningState if already into keepaliveState
		ap.setState(RunningState)
	case KeepAliveCmd:
		ap.sendAck(KeepAliveCmd)
	case StartCmd:
		ap.setState(StartingState)
		ap.eventListener.OnStarting(StartCmdCause)
	case StopCmd:
		ap.setState(StartingState)
		ap.eventListener.OnStarting(StopCmdCause)
		ap.sendAck(StopCmd)
	case GetCmd: // Get the state of an item and send it back.
		if length < 2 {
			ap.sendNack(GetCmd)
			return
		}
		itemId := ItemId(bytes[1])
		localItem := ap.GetLocalItemById(itemId)
		if localItem == nil {
			ap.sendNack(GetCmd)
			return
		}
		resultBuffer := make([]uint8, MaxOngetfuncBufferSize)
		resultBuffer = localItem.OnGetCmd(resultBuffer)
		resultLength := len(resultBuffer)
		if resultLength < 1 {
			ap.sendNack(GetCmd)
			return
		}
		if resultLength > MaxOngetfuncBufferSize { //TODO cant happen. Because if localItem.OnGetCmd() write outside the scope. something bad happened anyway.
			logMsg("Can't send full data. MAX_ONGETFUNC_BUFFER_SIZE is too small")
			resultLength = MaxOngetfuncBufferSize
		}
		ap.sendGetValue(itemId, resultBuffer)
	case GetValueCmd: // Response of "GET_CMD"
		if length < 3 {
			ap.sendNack(GetValueCmd)
			return
		}
		remoteItem := ap.GetRemoteItemById(ItemId(bytes[1]))
		if remoteItem == nil {
			ap.sendNack(GetValueCmd)
			return
		}
		remoteItem.OnGetValueCmd(bytes[2:])
		ap.sendAck(GetValueCmd)
		// reset timer to not fall into keepaliveState
		// or comeback into runningState if already into keepaliveState
		ap.setState(RunningState)
	case SetCmd:
		if length < 3 {
			ap.sendNack(SetCmd)
			return
		}
		localItem := ap.GetLocalItemById(ItemId(bytes[1]))
		if localItem == nil {
			ap.sendNack(SetCmd)
			return
		}
		localItem.OnSetCmd(bytes[2:])
		ap.sendAck(SetCmd)
	case StartStreamCmd:
		if length < 2 {
			ap.sendNack(StartCmd)
			return
		}
		result := ap.AddStreamSubscription(ItemId(bytes[1]))
		if !result {
			ap.sendNack(StartCmd)
			return
		}
		ap.sendAck(StartCmd)
	case StopStreamCmd:
		if length < 2 {
			ap.sendNack(StopStreamCmd)
			return
		}
		res := ap.RemoveStreamSubscription(ItemId(bytes[1]))
		if !res {
			ap.sendNack(StopStreamCmd)
			return
		}
		ap.sendAck(StopStreamCmd)
	case StreamValueCmd:
		if length < 3 {
			ap.sendNack(StreamValueCmd)
			return
		}
		remoteItem := ap.GetRemoteItemById(ItemId(bytes[1]))
		if remoteItem == nil {
			ap.sendNack(StreamValueCmd)
			return
		}
		remoteItem.OnGetValueCmd(bytes[2:])
		ap.sendAck(StreamValueCmd)
	default:
		ap.sendNack(command)
	}
}
