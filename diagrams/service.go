package diagrams

import (
	"errors"
	"sync"
	"systems-diagrams/systems"
)

var (
	ErrElementNotFound       = errors.New("element not found")
	ErrCantEditLinkedElement = errors.New("can't edit linked element")
)

type DiagramService struct {
	systemUpdatesChannel <-chan systems.System
	DiagramElements      map[int]*DiagramElement
	mu                   sync.Mutex
}

func NewDiagramService(systemUpdates <-chan systems.System) *DiagramService {
	return &DiagramService{
		systemUpdatesChannel: systemUpdates,
		DiagramElements:      make(map[int]*DiagramElement),
	}
}

func (s *DiagramService) ListenForUpdates() {
	for sys := range s.systemUpdatesChannel {
		s.mu.Lock()
		for _, element := range s.DiagramElements {
			if element.LinkedSystemID != nil && *element.LinkedSystemID == sys.Id {
				element.Name = sys.Name
				element.Description = sys.Description
			}
		}
		s.mu.Unlock()
	}
}

func (s *DiagramService) AddElement(element DiagramElement) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.DiagramElements[element.Id] = &element
}

func (s *DiagramService) EditElement(id int, name string, description string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	element, ok := s.DiagramElements[id]
	if !ok {
		return ErrElementNotFound
	}
	if element.LinkedSystemID != nil {
		return ErrCantEditLinkedElement
	}
	element.Name = name
	element.Description = description
	s.DiagramElements[id] = element
	return nil
}

func (s *DiagramService) UnlinkElement(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	element, ok := s.DiagramElements[id]
	if !ok {
		return ErrElementNotFound
	}
	element.LinkedSystemID = nil
	s.DiagramElements[id] = element
	return nil
}

func (s *DiagramService) LinkElement(id int, systemID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	element, ok := s.DiagramElements[id]
	if !ok {
		return ErrElementNotFound
	}
	element.LinkedSystemID = &systemID
	s.DiagramElements[id] = element
	return nil
}
