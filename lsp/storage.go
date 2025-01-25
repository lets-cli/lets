package lsp

type storage struct {
	documents map[string]*string
}

func newStorage() *storage {
	return &storage{
		documents: make(map[string]*string),
	}
}

func (s *storage) GetDocument(uri string) *string {
	return s.documents[uri]
}

func (s *storage) AddDocument(uri string, text string) {
	s.documents[uri] = &text
}