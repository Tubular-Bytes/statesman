package memory

import (
	"sync"

	"github.com/Tubular-Bytes/statesman/pkg/model"
)

type Storeable interface {
	model.LockData | model.State
}

type ItemStore[T Storeable] struct {
	mx    *sync.Mutex
	items map[string]T
}

func NewItemStore[T Storeable]() *ItemStore[T] {
	return &ItemStore[T]{
		mx:    &sync.Mutex{},
		items: make(map[string]T),
	}
}

func (s *ItemStore[T]) Get(id string) (T, bool) {
	s.mx.Lock()
	defer s.mx.Unlock()

	item, ok := s.items[id]

	return item, ok
}

func (s *ItemStore[T]) Put(id string, item T) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.items[id] = item
}

func (s *ItemStore[T]) Delete(id string) {
	s.mx.Lock()
	defer s.mx.Unlock()

	delete(s.items, id)
}

type Store struct {
	locks  *ItemStore[model.LockData]
	states *ItemStore[model.State]
}

func NewStore() *Store {
	return &Store{
		locks:  NewItemStore[model.LockData](),
		states: NewItemStore[model.State](),
	}
}

/**
	GetState(id string) (*model.State, error)
	PutState(id string, state *model.State) error

	Lock(lockData *model.LockData) error
	Unlock(lockID string) error
**/

func (s *Store) GetState(id string) (*model.State, error) {
	if state, ok := s.states.Get(id); ok {
		return &state, nil
	}

	return nil, model.ErrNotFound
}

func (s *Store) PutState(id string, state *model.State) error {
	if state == nil {
		return model.ErrInvalidState
	}

	s.states.Put(id, *state)

	return nil
}

func (s *Store) Lock(lockData *model.LockData) error {
	if lockData == nil {
		return model.ErrInvalidLock
	}

	if _, ok := s.locks.Get(lockData.LockID); ok {
		return model.ErrLockConflict
	}

	s.locks.Put(lockData.LockID, *lockData)

	return nil
}

func (s *Store) Unlock(lockID string) error {
	if _, ok := s.locks.Get(lockID); !ok {
		return model.ErrNotFound
	}

	s.locks.Delete(lockID)

	return nil
}
