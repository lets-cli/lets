package lsp

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMixinsStoresAndIndexesMixinDocuments(t *testing.T) {
	dir := t.TempDir()

	mainPath := filepath.Join(dir, "lets.yaml")
	baseMixinPath := filepath.Join(dir, "lets.base.yaml")
	localMixinPath := filepath.Join(dir, "lets.local.yaml")

	baseMixinDoc := `commands:
  build:
    cmd: echo build`

	localMixinDoc := `commands:
  test:
    cmd: echo test`

	if err := os.WriteFile(baseMixinPath, []byte(baseMixinDoc), 0o644); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", baseMixinPath, err)
	}

	if err := os.WriteFile(localMixinPath, []byte(localMixinDoc), 0o644); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", localMixinPath, err)
	}

	mainDoc := `mixins:
  - lets.base.yaml
  - -lets.local.yaml
commands:
  release:
    depends: [build, test]
    cmd: echo release`

	server := &lspServer{
		storage: newStorage(),
		parser:  newParser(logger),
		index:   newIndex(logger),
		log:     logger,
	}

	mainURI := pathToURI(mainPath)
	server.storage.AddDocument(mainURI, mainDoc)
	server.loadMixins(mainURI)

	baseMixinURI := pathToURI(baseMixinPath)
	localMixinURI := pathToURI(localMixinPath)

	if got := server.storage.GetDocument(baseMixinURI); got == nil || *got != baseMixinDoc {
		t.Fatalf("storage for %s = %#v, want %q", baseMixinURI, got, baseMixinDoc)
	}

	if got := server.storage.GetDocument(localMixinURI); got == nil || *got != localMixinDoc {
		t.Fatalf("storage for %s = %#v, want %q", localMixinURI, got, localMixinDoc)
	}

	buildInfo, ok := server.index.findCommand("build")
	if !ok {
		t.Fatal("expected build command from mixin to be indexed")
	}

	if buildInfo.fileURI != baseMixinURI {
		t.Fatalf("build indexed at %s, want %s", buildInfo.fileURI, baseMixinURI)
	}

	testInfo, ok := server.index.findCommand("test")
	if !ok {
		t.Fatal("expected test command from mixin to be indexed")
	}

	if testInfo.fileURI != localMixinURI {
		t.Fatalf("test indexed at %s, want %s", testInfo.fileURI, localMixinURI)
	}
}

func TestLoadMixinsSkipsMissingFiles(t *testing.T) {
	dir := t.TempDir()

	mainPath := filepath.Join(dir, "lets.yaml")
	mainURI := pathToURI(mainPath)

	mainDoc := `mixins:
  - missing.yaml
commands:
  release:
    cmd: echo release`

	server := &lspServer{
		storage: newStorage(),
		parser:  newParser(logger),
		index:   newIndex(logger),
		log:     logger,
	}

	server.storage.AddDocument(mainURI, mainDoc)
	server.loadMixins(mainURI)

	if got := server.storage.GetDocument(pathToURI(filepath.Join(dir, "missing.yaml"))); got != nil {
		t.Fatalf("expected missing mixin to not be stored, got %#v", got)
	}

	if _, ok := server.index.findCommand("missing"); ok {
		t.Fatal("expected no indexed command for missing mixin")
	}
}
