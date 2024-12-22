package local

import (
	"sync"

	"code.emcdtech.com/b2b/swap/model"
	"github.com/google/uuid"
)

type Subscribers struct {
	mu           sync.RWMutex
	subsBySwapID map[uuid.UUID][]*model.Subscriber
}

func NewSubscribers() *Subscribers {
	return &Subscribers{
		subsBySwapID: make(map[uuid.UUID][]*model.Subscriber),
	}
}

func (s *Subscribers) Add(swapID, clientID uuid.UUID, ch chan model.PublicStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()

	subs, ok := s.subsBySwapID[swapID]
	if !ok {
		subs = []*model.Subscriber{}
		s.subsBySwapID[swapID] = subs
	}

	s.subsBySwapID[swapID] = append(subs, &model.Subscriber{ClientID: clientID, Ch: ch})
}

func (s *Subscribers) GetBySwapID(swapID uuid.UUID) []*model.Subscriber {
	s.mu.RLock()
	defer s.mu.RUnlock()

	subs, ok := s.subsBySwapID[swapID]
	if !ok {
		return nil
	}

	// Копируем указатели, что бы избежать проблем при итерации по слайсу,
	// при изменении длины оригинального слайса
	subsCopy := make([]*model.Subscriber, len(subs))
	copy(subsCopy, subs)

	return subsCopy
}

// Delete удаляет клиента.
// Не закрывает канал, что бы не допустить паники в broadcast на уровне service
func (s *Subscribers) Delete(swapID, clientID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	subs, ok := s.subsBySwapID[swapID]
	if !ok {
		return
	}

	var newSubs []*model.Subscriber

	for _, sub := range subs {
		if sub.ClientID != clientID {
			newSubs = append(newSubs, sub)
		}
	}

	if len(newSubs) == 0 {
		delete(s.subsBySwapID, swapID)
	} else {
		s.subsBySwapID[swapID] = newSubs
	}
}
