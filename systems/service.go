package systems

import (
	"errors"
	"sync"
)

var (
	ErrSystemNotFound      = errors.New("system not found")
	ErrSystemAlreadyExists = errors.New("system already exists")
)

type SystemsService struct {
	Systems           map[int]*System
	systemUpdatesChan chan<- System
	mu                sync.Mutex
}

func NewSystemsService(systemUpdates chan<- System) *SystemsService {
	return &SystemsService{
		Systems:           make(map[int]*System),
		systemUpdatesChan: systemUpdates,
	}
}

func (s *SystemsService) AddSystem(system System) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Systems[system.Id]; ok {
		return ErrSystemAlreadyExists
	}
	s.Systems[system.Id] = &system
	return nil
}

func (s *SystemsService) UpdateSystem(id int, name string, description string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	system, ok := s.Systems[id]
	if !ok {
		return ErrSystemNotFound
	}
	system.Name = name
	system.Description = description
	s.Systems[id] = system
	s.systemUpdatesChan <- *system
	return nil
}

func (s *SystemsService) GetSystem(id int) (*System, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	system, ok := s.Systems[id]
	if !ok {
		return nil, ErrSystemNotFound
	}
	return system, nil
}

func (s *SystemsService) DeleteSystem(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Systems[id]; !ok {
		return ErrSystemNotFound
	}
	delete(s.Systems, id)
	return nil
}
