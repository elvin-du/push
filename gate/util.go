package main

import (
	"push/common/model"
)

func ValidateRegID(id string) error {
	_, err := model.RegistryModel().Get(id)
	if nil != err {
		return err
	}

	return nil
}
