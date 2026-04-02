package lsp

import (
	"reflect"
	"testing"
)

func TestIndexDocumentStoresCommands(t *testing.T) {
	doc := `commands:
  build:
    cmd: echo build
  test:
    cmd: echo test`

	idx := newIndex(logger)
	idx.IndexDocument("file:///tmp/lets.yaml", doc)

	tests := []struct {
		name string
		want commandInfo
	}{
		{
			name: "build",
			want: commandInfo{
				fileURI:  "file:///tmp/lets.yaml",
				position: pos(1, 2),
			},
		},
		{
			name: "test",
			want: commandInfo{
				fileURI:  "file:///tmp/lets.yaml",
				position: pos(3, 2),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := idx.findCommand(tt.name)
			if !ok {
				t.Fatalf("findCommand(%q) did not find command", tt.name)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("findCommand(%q) = %#v, want %#v", tt.name, got, tt.want)
			}
		})
	}
}

func TestIndexDocumentReplacesCommandsForSameDocument(t *testing.T) {
	originalDoc := `commands:
  build:
    cmd: echo build
  test:
    cmd: echo test`

	updatedDoc := `commands:
  release:
    depends: [build]
    cmd: echo release`

	idx := newIndex(logger)
	idx.IndexDocument("file:///tmp/lets.yaml", originalDoc)
	idx.IndexDocument("file:///tmp/lets.yaml", updatedDoc)

	if _, ok := idx.findCommand("build"); ok {
		t.Fatal("expected build to be removed after reindex")
	}

	if _, ok := idx.findCommand("test"); ok {
		t.Fatal("expected test to be removed after reindex")
	}

	got, ok := idx.findCommand("release")
	if !ok {
		t.Fatal("expected release to be indexed after reindex")
	}

	want := commandInfo{
		fileURI:  "file:///tmp/lets.yaml",
		position: pos(1, 2),
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("findCommand(%q) = %#v, want %#v", "release", got, want)
	}
}
