package applayerprotocol

type SimpleLocalItem struct {
	Id                 ItemId
	OnSetCmdCallback   func(bytes []uint8)
	OnGetCmdCallback   func(bytes []uint8) []uint8
	HasChangedCallback func() bool
}

func (sli *SimpleLocalItem) GetId() ItemId {
	return sli.Id
}
func (sli *SimpleLocalItem) OnSetCmd(bytes []uint8) {
	if sli.OnSetCmdCallback == nil {
		return
	}
	sli.OnSetCmdCallback(bytes)
}
func (sli *SimpleLocalItem) OnGetCmd(bytes []uint8) []uint8 {
	if sli.OnGetCmdCallback == nil {
		return nil
	}
	return sli.OnGetCmdCallback(bytes)
}
func (sli *SimpleLocalItem) HasChanged() bool {
	if sli.HasChangedCallback == nil {
		return false
	}
	return sli.HasChangedCallback()
}
