package dal

type Mysql struct {
}

func (*Mysql) GetUserInfoByClientId(cliId string) error {
	return nil
}

func (*Mysql) SaveOfflineMsg() error {
	return nil
}

func (*Mysql) GetGateInfoByUserId(userId string) (*GateInfo, error) {
	return &GateInfo{}, nil
}
