package model

type OfflineMsg struct {
	ClientID string
	PacketID uint16
	Kind     int
	Content  string
	Extra    string
	CreateAt uint64
}
