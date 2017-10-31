package model

var (
	_apps []*App = nil
)

func GetApps() ([]*App, error) {
	if nil == _apps {
		var err error = nil
		_apps, err = AppModel().GetAll()
		if nil != err {
			return nil, err
		}
	}

	return _apps, nil
}
