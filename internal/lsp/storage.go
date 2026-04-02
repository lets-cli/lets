package lsp

import "sync"

type storage struct {
	mu        sync.RWMutex
	documents map[string]*string
}

func newStorage() *storage {
	return &storage{
		documents: make(map[string]*string),
	}
}

func (s *storage) GetDocument(uri string) *string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.documents[uri]
}

func (s *storage) AddDocument(uri string, text string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.documents[uri] = &text
}
