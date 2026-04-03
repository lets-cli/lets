package lsp

import "sync"

type storage struct {
	mu        sync.RWMutex
	documents map[string]*string
	mixins    map[string][]string
}

func newStorage() *storage {
	return &storage{
		documents: make(map[string]*string),
		mixins:    make(map[string][]string),
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

func (s *storage) SetMixins(uri string, mixins []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.mixins[uri] = mixins
}

func (s *storage) GetMixins(uri string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.mixins[uri]
}
