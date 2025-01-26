package lsp

import (
	"fmt"
	"slices"

	"github.com/lets-cli/lets/util"
	"github.com/tliron/glsp"
	lsp "github.com/tliron/glsp/protocol_3_16"
)

func (s *lspServer) initialize(context *glsp.Context, params *lsp.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()
	syncKind := lsp.TextDocumentSyncKindFull
	capabilities.TextDocumentSync.(*lsp.TextDocumentSyncOptions).Change = &syncKind

	capabilities.CompletionProvider = &lsp.CompletionOptions{
		TriggerCharacters: []string{" ", "- ", "["},
	}

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

type definitionHandler struct {
	parser *parser
}

func (h *definitionHandler) findMixinsDefinition(doc *string, params *lsp.DefinitionParams) (any, error) {
	path := normalizePath(params.TextDocument.URI)
	filename := h.parser.extractFilenameFromMixins(doc, params.Position)
	if filename == "" {
		return nil, nil
	}

	absFilename := replacePathFilename(path, filename)

	if !util.FileExists(absFilename) {
		return nil, nil
	}

	return []lsp.Location{
		{
			URI:   pathToUri(absFilename),
			Range: lsp.Range{},
		},
	}, nil
}

type completionHandler struct {
	parser *parser
}

func (h *completionHandler) buildDependsCompletions(doc *string, params *lsp.CompletionParams) ([]lsp.CompletionItem, error) {
	commands := h.parser.getCommands(doc)
	alreadyInDepends := h.parser.extractDependsValues(doc)
	currentCommand := h.parser.getCurrentCommand(doc, params.Position)
	items := []lsp.CompletionItem{}

	keywordKind := lsp.CompletionItemKindKeyword

	for _, cmd := range commands {
		// do not suggest the current command
		if currentCommand != nil && cmd.name == currentCommand.name {
			continue
		}
		// do not suggest already included commands
		if slices.Contains(alreadyInDepends, cmd.name) {
			continue
		}
		items = append(items, lsp.CompletionItem{
			Label: cmd.name,
			Kind:  &keywordKind,
		})
	}

	return items, nil
}

// Returns: Location | []Location | []LocationLink | nil
func (s *lspServer) textDocumentDefinition(context *glsp.Context, params *lsp.DefinitionParams) (any, error) {
	definitionHandler := definitionHandler{
		parser: newParser(s.log),
	}
	doc := s.storage.GetDocument(params.TextDocument.URI)

	p := newParser(s.log)
	switch p.getPositionType(doc, params.Position) {
	case PositionTypeMixins:
		return definitionHandler.findMixinsDefinition(doc, params)
	default:
		return nil, nil
	}
}

// Returns: []CompletionItem | CompletionList | nil
func (s *lspServer) textDocumentCompletion(context *glsp.Context, params *lsp.CompletionParams) (any, error) {
	completionHandler := completionHandler{
		parser: newParser(s.log),
	}
	doc := s.storage.GetDocument(params.TextDocument.URI)

	p := newParser(s.log)
	pos := p.getPositionType(doc, params.Position)
	s.log.Infof("Position type: %d", pos)
	switch pos {
	case PositionTypeDepends:
		return completionHandler.buildDependsCompletions(doc, params)
	default:
		return []lsp.CompletionItem{}, nil
	}
}
