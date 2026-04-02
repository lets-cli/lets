package lsp

import (
	"errors"
	"slices"

	"github.com/lets-cli/lets/internal/util"
	"github.com/tliron/commonlog"
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
	go s.index.IndexDocument(params.TextDocument.URI, params.TextDocument.Text)
	return nil
}

func (s *lspServer) textDocumentDidChange(context *glsp.Context, params *lsp.DidChangeTextDocumentParams) error {
	for _, change := range params.ContentChanges {
		switch c := change.(type) {
		case lsp.TextDocumentContentChangeEventWhole:
			s.storage.AddDocument(params.TextDocument.URI, c.Text)
			go s.index.IndexDocument(params.TextDocument.URI, c.Text)
		case lsp.TextDocumentContentChangeEvent:
			return errors.New("incremental changes not supported")
		}
	}

	return nil
}

type definitionHandler struct {
	log    commonlog.Logger
	parser *parser
	index  *index
}

func (h *definitionHandler) findMixinsDefinition(doc *string, params *lsp.DefinitionParams) (any, error) {
	path := normalizePath(params.TextDocument.URI)

	filename := h.parser.extractFilenameFromMixins(doc, params.Position)
	if filename == "" {
		h.parser.log.Debugf("no mixin filename resolved at %s:%d:%d", path, params.Position.Line, params.Position.Character)
		return nil, nil
	}

	absFilename := replacePathFilename(path, filename)

	if !util.FileExists(absFilename) {
		h.parser.log.Debugf("mixin target does not exist: %s", absFilename)
		return nil, nil
	}

	h.parser.log.Debugf("resolved mixin definition %q -> %s", filename, absFilename)

	return []lsp.Location{
		{
			URI:   pathToURI(absFilename),
			Range: lsp.Range{},
		},
	}, nil
}

func locationForCommand(uri string, position lsp.Position) lsp.Location {
	return lsp.Location{
		URI: uri,
		Range: lsp.Range{
			Start: lsp.Position{
				Line:      position.Line,
				Character: 2, // TODO: do we have to assume indentation?
			},
			End: lsp.Position{
				Line:      position.Line,
				Character: 2, // TODO: do we need + len ?
			},
		},
	}
}

func (h *definitionHandler) findCommandDefinition(doc *string, params *lsp.DefinitionParams) (any, error) {
	path := normalizePath(params.TextDocument.URI)

	commandName := h.parser.extractCommandReference(doc, params.Position)
	if commandName == "" {
		h.log.Debugf("no command reference resolved at %s:%d:%d", path, params.Position.Line, params.Position.Character)
		return nil, nil
	}

	commandInfo, found := h.index.findCommand(commandName)
	if !found {
		h.log.Debugf("command reference %q did not match any local command", commandName)
		return nil, nil
	}

	h.log.Debugf(
		"resolved command definition %q -> %s:%d:%d",
		commandName,
		path,
		commandInfo.position.Line,
		commandInfo.position.Character,
	)

	loc := locationForCommand(commandInfo.fileURI, commandInfo.position)
	return []lsp.Location{loc}, nil
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

// Returns: Location | []Location | []LocationLink | nil.
func (s *lspServer) textDocumentDefinition(context *glsp.Context, params *lsp.DefinitionParams) (any, error) {
	definitionHandler := definitionHandler{
		log:    s.log,
		parser: newParser(s.log),
		index:  s.index,
	}
	doc := s.storage.GetDocument(params.TextDocument.URI)

	p := newParser(s.log)
	positionType := p.getPositionType(doc, params.Position)
	s.log.Debugf(
		"definition request uri=%s line=%d char=%d type=%s",
		normalizePath(params.TextDocument.URI),
		params.Position.Line,
		params.Position.Character,
		positionType,
	)

	switch positionType {
	case PositionTypeMixins:
		return definitionHandler.findMixinsDefinition(doc, params)
	case PositionTypeDepends, PositionTypeCommandAlias:
		return definitionHandler.findCommandDefinition(doc, params)
	default:
		s.log.Debugf("definition request ignored: unsupported cursor position")
		return nil, nil
	}
}

// Returns: []CompletionItem | CompletionList | nil.
func (s *lspServer) textDocumentCompletion(context *glsp.Context, params *lsp.CompletionParams) (any, error) {
	completionHandler := completionHandler{
		parser: newParser(s.log),
	}
	doc := s.storage.GetDocument(params.TextDocument.URI)

	p := newParser(s.log)
	switch p.getPositionType(doc, params.Position) {
	case PositionTypeDepends:
		return completionHandler.buildDependsCompletions(doc, params)
	default:
		return []lsp.CompletionItem{}, nil
	}
}
