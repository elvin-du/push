package dal

type Redis struct {
}

func (*Redis) GetUserInfoByClientId(cliId string) error {
	return nil
}

func (*Redis) SaveOfflineMsg() error {
	return nil
}

func (*Redis) GetGateInfoByUserId(userId string) (*GateInfo, error) {
	return &GateInfo{}, nil
}
