package lsp

import (
	"context"
	"fmt"

	"github.com/tliron/commonlog"
	_ "github.com/tliron/commonlog/simple"
	lsp "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

const lsName = "lets_ls"

var (
	handler lsp.Handler
)

type lspServer struct {
	version string
	server  *server.Server
	storage *storage
	log     commonlog.Logger
}

func (s *lspServer) Run() error {
	return s.server.RunStdio()
}

func Run(ctx context.Context, version string) error {
	commonlog.Configure(1, nil)
	logger := commonlog.GetLogger(fmt.Sprintf("%s.parser", lsName))
	logger.Info("Lets LSP server starting")

	handler = lsp.Handler{}

	glspServer := server.NewServer(&handler, lsName, false)
	glspServer.Context = ctx

	lspServer := &lspServer{
		version: version,
		server:  glspServer,
		storage: newStorage(),
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

	return lspServer.Run()
}
