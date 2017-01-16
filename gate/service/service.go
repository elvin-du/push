package service

import (
	"sync"
)

type Processor interface {
	Process() error
}

type Receiver interface {
	ReadLoop() error
}

type Sendor interface {
	WriteLoop() error
}

type Service struct {
	clientId string
	Processor
	Sendor
	Receiver
	wgStarted sync.WaitGroup
}

func NewService(cliId string) *Service {
	return &Service{
		clientId: cliId,
	}
}

func (s *Service) SetId(id string) *Service {
	s.clientId = id
	return s
}

func (s *Service) SetProcessor(p Processor) *Service {
	s.Processor = p
	return s
}

func (s *Service) SetSendor(sd Sendor) *Service {
	s.Sendor = sd
	return s
}

func (s *Service) SetReceiver(r Receiver) *Service {
	s.Receiver = r
	return s
}

func (s *Service) Start() error {
	s.wgStarted.Add(1)
	go s.Process()

	s.wgStarted.Add(1)
	go s.ReadLoop()

	s.wgStarted.Add(1)
	go s.WriteLoop()

	s.wgStarted.Wait()

	return nil
}
