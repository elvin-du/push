package mqtt

import (
	"github.com/surgemq/message"
)

func (s *Service) Process(msg message.Message) error {
	var err error = nil

	switch msg := msg.(type) {
	case *message.PublishMessage:
		err = s.processPublish(msg)
	case *message.PubackMessage:
	case *message.DisconnectMessage:
	case *message.PingreqMessage:
	case *message.PingrespMessage:
	default:
	}

	return err
}

func (s *Service) processPublish(msg *message.PublishMessage) error {
	return nil
}
