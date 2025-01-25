package lsp

import (
	// "log/slog"

	"context"
	"fmt"

	"github.com/tliron/commonlog"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
	log "github.com/sirupsen/logrus"

	// Must include a backend implementation
	// See CommonLog for other options: https://github.com/tliron/commonlog
	_ "github.com/tliron/commonlog/simple"
)

const lsName = "lets_ls"

var (
	handler protocol.Handler
)

type lspServer struct {
	version string
	server *server.Server
	storage *storage
}

func (s *lspServer) Run() error {
	return s.server.RunStdio()
}

func Run(ctx context.Context, version string) error {
	fmt.Println("lets: LSP server starting fmt")
	commonlog.Configure(1, nil)
	logger := commonlog.GetLogger(fmt.Sprintf("%s.parser", lsName))
	logger.Info("lets: LSP server starting")
	log.Info("lets: LSP server starting logrst")

	handler = protocol.Handler{}

	glspServer := server.NewServer(&handler, lsName, false)
	glspServer.Context = ctx

	lspServer := &lspServer{
		version: version,
		server: glspServer,
		storage: newStorage(),
	}

	handler.Initialize = lspServer.initialize
	handler.Initialized = lspServer.initialized
	handler.Shutdown = lspServer.shutdown
	handler.SetTrace = lspServer.setTrace
	handler.TextDocumentDidOpen = lspServer.textDocumentDidOpen
	handler.TextDocumentDidChange = lspServer.textDocumentDidChange
	handler.TextDocumentDefinition = lspServer.textDocumentDefinition

	return lspServer.Run()
}
