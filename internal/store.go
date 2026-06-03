package handlers

import (
	"log"
	"sync"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]interface{}),
	}
}

func (s *Store) Set(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *Store) Get(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, exists := s.data[key]
	return val, exists
}

func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.data[key]; !exists {
		log.Fatal("Key not found")
	}
	delete(s.data, key)
}
