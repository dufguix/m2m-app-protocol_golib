package applayerprotocol

type SimpleRemoteItem struct {
	Id                    ItemId
	OnGetValueCmdCallback func([]uint8)
}

func (sri *SimpleRemoteItem) GetId() ItemId {
	return sri.Id
}
func (sri *SimpleRemoteItem) OnGetValueCmd(bytes []uint8) {
	if sri.OnGetValueCmdCallback == nil {
		return
	}
	sri.OnGetValueCmdCallback(bytes)
}
