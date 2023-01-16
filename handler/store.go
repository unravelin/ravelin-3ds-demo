package handler

import (
	"errors"
	"sync"
)

const (
	MethodStatusCompleted    = "Y"
	MethodStatusNotCompleted = "N"
	MethodStatusUnavailable  = "U"
)

type ThreeDSTransaction struct {
	MessageVersion string
	MethodStatus   string
}

type ThreeDSTransactionStore struct {
	mu    *sync.RWMutex
	store map[string]ThreeDSTransaction
}

func NewThreeDSTransactionStore() ThreeDSTransactionStore {
	return ThreeDSTransactionStore{
		mu:    &sync.RWMutex{},
		store: make(map[string]ThreeDSTransaction),
	}
}

func (s *ThreeDSTransactionStore) Add(threeDSServerTransID string, tx ThreeDSTransaction) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[threeDSServerTransID] = tx
}

func (s *ThreeDSTransactionStore) Get(threeDSTransactionID string) (ThreeDSTransaction, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tx, ok := s.store[threeDSTransactionID]
	return tx, ok
}

func (s *ThreeDSTransactionStore) SetMethodStatus(threeDSTransactionID string, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, ok := s.store[threeDSTransactionID]
	if !ok {
		return errors.New("not found")
	}

	tx.MethodStatus = status
	s.store[threeDSTransactionID] = tx
	return nil
}
