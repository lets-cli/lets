package lsp

import (
	"maps"
	"sync"

	"github.com/lets-cli/lets/internal/set"
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
	parser        *parser
	mu            sync.RWMutex
	commands      map[string]commandInfo
	commandsByURI map[string]set.Set[string]
}

func newIndex(log commonlog.Logger) *index {
	return &index{
		log:           log,
		parser:        newParser(log),
		commands:      make(map[string]commandInfo),
		commandsByURI: make(map[string]set.Set[string]),
	}
}

// IndexDocument extracts commands from a document and updates the index to reflect that document's current state.
func (i *index) IndexDocument(uri string, doc string) {
	commands := i.parser.getCommands(&doc)

	indexedCommands := make(map[string]commandInfo, len(commands))
	indexedNames := set.NewSet[string]()

	for _, command := range commands {
		indexedCommands[command.name] = commandInfo{
			fileURI:  uri,
			position: command.position,
		}
		indexedNames.Add(command.name)
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	i.log.Debugf("Indexed %d commands in file %s", len(indexedNames), uri)

	for name := range i.commandsByURI[uri] {
		delete(i.commands, name)
	}

	maps.Copy(i.commands, indexedCommands)

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
