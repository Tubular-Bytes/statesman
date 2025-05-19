package model

import "fmt"

var (
	ErrNotFound     = fmt.Errorf("not found")
	ErrInvalidState = fmt.Errorf("invalid state")
	ErrInvalidLock  = fmt.Errorf("invalid lock")
	ErrLockConflict = fmt.Errorf("lock conflict")
)
