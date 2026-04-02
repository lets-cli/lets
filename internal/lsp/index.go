package lsp

import (
	"sync"

	"github.com/tliron/commonlog"
	lsp "github.com/tliron/glsp/protocol_3_16"
)

// TODO: maybe use Command struct ?
type commandInfo struct {
	fileURI string
	// position stored at the time of indexing and may be stale
	position lsp.Position
}

type index struct {
	log           commonlog.Logger
	mu            sync.RWMutex
	commands      map[string]commandInfo
	commandsByURI map[string]map[string]struct{}
}

func newIndex(log commonlog.Logger) *index {
	return &index{
		log:           log,
		commands:      make(map[string]commandInfo),
		commandsByURI: make(map[string]map[string]struct{}),
	}
}

// IndexDocument extracts commands from a document and updates the index to reflect that document's current state.
func (i *index) IndexDocument(uri string, doc string) {
	parser := newParser(i.log)
	commands := parser.getCommands(&doc)

	indexedCommands := make(map[string]commandInfo, len(commands))
	indexedNames := make(map[string]struct{}, len(commands))

	for _, command := range commands {
		indexedCommands[command.name] = commandInfo{
			fileURI:  uri,
			position: command.position,
		}
		// TODOL maybe use Set
		indexedNames[command.name] = struct{}{}
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	i.log.Debugf("Indexed %d commands in file %s", len(indexedNames), uri)

	for name := range i.commandsByURI[uri] {
		delete(i.commands, name)
	}

	for name, info := range indexedCommands {
		i.commands[name] = info
	}

	if len(indexedNames) == 0 {
		delete(i.commandsByURI, uri)
		return
	}

	i.commandsByURI[uri] = indexedNames
}

func (i *index) findCommand(name string) (commandInfo, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	command, ok := i.commands[name]
	return command, ok
}
