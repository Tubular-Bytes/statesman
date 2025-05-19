package backend

import (
	"github.com/Tubular-Bytes/statesman/pkg/backend/memory"
	"github.com/Tubular-Bytes/statesman/pkg/model"
)

type Store interface {
	GetState(id string) (*model.State, error)
	PutState(id string, state *model.State) error

	Lock(lockData *model.LockData) error
	Unlock(lockID string) error
}

var _ Store = (*memory.Store)(nil)

var store Store

func Get() Store {
	if store == nil {
		initStore()
	}
	return store
}

func initStore() {
	store = memory.NewStore()
}
