package lsp

import (
	"context"

	"github.com/lets-cli/lets/internal/env"
	"github.com/tliron/commonlog"
	_ "github.com/tliron/commonlog/simple"
	lsp "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

const lsName = "lets_ls"

var handler lsp.Handler

type lspServer struct {
	version string
	server  *server.Server
	storage *storage
	index   *index
	log     commonlog.Logger
}

func (s *lspServer) Run() error {
	return s.server.RunStdio()
}

func lspLogVerbosity() int {
	verbosity := 1

	defer func() {
		_ = recover()
	}()

	if env.DebugLevel() > 0 {
		verbosity = 2
	}

	return verbosity
}

func Run(ctx context.Context, version string) error {
	commonlog.Configure(lspLogVerbosity(), nil)

	logger := commonlog.GetLogger(lsName)
	logger.Infof("Lets LSP server starting %s", version)

	handler = lsp.Handler{}

	glspServer := server.NewServer(&handler, lsName, false)
	glspServer.Context = ctx

	lspServer := &lspServer{
		version: version,
		server:  glspServer,
		storage: newStorage(),
		index:   newIndex(logger),
		log:     logger,
	}

	handler.Initialize = lspServer.initialize
	handler.Initialized = lspServer.initialized
	handler.Shutdown = lspServer.shutdown
	handler.SetTrace = lspServer.setTrace
	handler.TextDocumentDidOpen = lspServer.textDocumentDidOpen
	handler.TextDocumentDidChange = lspServer.textDocumentDidChange
	handler.TextDocumentDefinition = lspServer.textDocumentDefinition
	handler.TextDocumentCompletion = lspServer.textDocumentCompletion
	// TODO: add onDelete

	return lspServer.Run()
}
