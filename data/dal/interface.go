package dal

type GateInfo struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

type DAL interface {
	GetUserInfoByClientId(cliId string) error

	SaveOfflineMsg() error

	GetGateInfoByUserId(userId string) (*GateInfo, error)
}
