package model

type OfflineMsg struct {
	AppID    string
	ClientID string
	PacketID uint16
	Kind     int
	Content  string
	Extra    string
	CreateAt uint64
}

type App struct {
	ID          string
	Secret      string
	AuthType    uint16
	Name        string
	Description string
	Status      byte
	CreateAt    uint64
	UpdatedAt   uint64
}
