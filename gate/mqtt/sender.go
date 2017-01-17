package mqtt

import (
	"errors"
)

var (
	E_READ_ERROR = errors.New("read error")
)

func (s *Service) ReadLoop() error {
	for {

	}

	return E_READ_ERROR
}
