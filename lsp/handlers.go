package lsp

import (
	"fmt"

	"github.com/lets-cli/lets/util"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *lspServer) initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &s.version,
		},
	}, nil
}

func (s *lspServer) initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func (s *lspServer) shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func (s *lspServer) setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func (s *lspServer) textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	s.storage.AddDocument(params.TextDocument.URI, params.TextDocument.Text)
	return nil
}

func (s *lspServer) textDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	for _, change := range params.ContentChanges {
		switch c := change.(type) {
			case protocol.TextDocumentContentChangeEventWhole:
				s.storage.AddDocument(params.TextDocument.URI, c.Text)
			case protocol.TextDocumentContentChangeEvent:
				return fmt.Errorf("incremental changes not supported")
		}
	}
	return nil
}

type DefinitionHandler struct {}

func (h *DefinitionHandler) findMixinsDefinition(doc *string, params *protocol.DefinitionParams) (any, error) {
	path := normalizePath(params.TextDocument.URI)
	filename := extractFilenameFromMixins(doc, params.Position)
	if filename == "" {
		return nil, nil
	}

	absFilename := replacePathFilename(path, filename)

	if !util.FileExists(absFilename) {
		return nil, nil
	}

	return []protocol.Location{
		{
			URI: pathToUri(absFilename),
			Range: protocol.Range{},
		},
	}, nil
}

// Returns: Location | []Location | []LocationLink | nil
func (s *lspServer) textDocumentDefinition(context *glsp.Context, params *protocol.DefinitionParams) (any, error) {
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
