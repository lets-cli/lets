package lsp

import (
	"fmt"

	"github.com/lets-cli/lets/util"
	"github.com/tliron/glsp"
	lsp "github.com/tliron/glsp/protocol_3_16"
)

func (s *lspServer) initialize(context *glsp.Context, params *lsp.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()
	value := lsp.TextDocumentSyncKindFull
	capabilities.TextDocumentSync.(*lsp.TextDocumentSyncOptions).Change = &value

	return lsp.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &lsp.InitializeResultServerInfo{
			Name:    lsName,
			Version: &s.version,
		},
	}, nil
}

func (s *lspServer) initialized(context *glsp.Context, params *lsp.InitializedParams) error {
	return nil
}

func (s *lspServer) shutdown(context *glsp.Context) error {
	lsp.SetTraceValue(lsp.TraceValueOff)
	return nil
}

func (s *lspServer) setTrace(context *glsp.Context, params *lsp.SetTraceParams) error {
	lsp.SetTraceValue(params.Value)
	return nil
}

func (s *lspServer) textDocumentDidOpen(context *glsp.Context, params *lsp.DidOpenTextDocumentParams) error {
	s.storage.AddDocument(params.TextDocument.URI, params.TextDocument.Text)
	return nil
}

func (s *lspServer) textDocumentDidChange(context *glsp.Context, params *lsp.DidChangeTextDocumentParams) error {
	for _, change := range params.ContentChanges {
		switch c := change.(type) {
			case lsp.TextDocumentContentChangeEventWhole:
				s.storage.AddDocument(params.TextDocument.URI, c.Text)
			case lsp.TextDocumentContentChangeEvent:
				return fmt.Errorf("incremental changes not supported")
		}
	}
	return nil
}

type DefinitionHandler struct {}

func (h *DefinitionHandler) findMixinsDefinition(doc *string, params *lsp.DefinitionParams) (any, error) {
	path := normalizePath(params.TextDocument.URI)
	filename := extractFilenameFromMixins(doc, params.Position)
	if filename == "" {
		return nil, nil
	}

	absFilename := replacePathFilename(path, filename)

	if !util.FileExists(absFilename) {
		return nil, nil
	}

	return []lsp.Location{
		{
			URI: pathToUri(absFilename),
			Range: lsp.Range{},
		},
	}, nil
}

// Returns: Location | []Location | []LocationLink | nil
func (s *lspServer) textDocumentDefinition(context *glsp.Context, params *lsp.DefinitionParams) (any, error) {
	definitionHandler := DefinitionHandler{}
	doc := s.storage.GetDocument(params.TextDocument.URI)

	switch getPositionType(doc, params.Position) {
		case PositionTypeMixins:
			return definitionHandler.findMixinsDefinition(doc, params)
		case PositionTypeNone:
			return nil, nil
	}
	return nil, nil
}
