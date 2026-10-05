package astparser

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestParseGraphqlDocumentFileImports(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "graphql-imports")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	mkdirAll := func(path string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Join(tempDir, path), 0755); err != nil {
			t.Fatal(err)
		}
	}
	writeFile := func(path, content string) {
		t.Helper()
		if err := ioutil.WriteFile(filepath.Join(tempDir, path), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	mkdirAll("app")
	mkdirAll("common/specific")
	mkdirAll("common/types")

	writeFile("app/main.graphql", "#import \"../common/schema.graphql\"\n\ntype Main { id: ID }\n")
	writeFile("common/schema.graphql", "#import \"./specific/schema.graphql\"\n#import \"./types/*.graphql\"\n#import \"./types/*.graphql\"\n\ntype Common { id: ID }\n")
	writeFile("common/specific/schema.graphql", "type Specific { id: ID }\n")
	writeFile("common/types/a.graphql", "type A { id: ID }\n")
	writeFile("common/types/b.graphql", "type B { id: ID }\n")
	writeFile("common/types/not-graphql.txt", "type Text { id: ID }\n")

	doc, report := ParseGraphqlDocumentFile(filepath.Join(tempDir, "app/main.graphql"))
	if report.HasErrors() {
		t.Fatal(report.Error())
	}

	var names []string
	for _, def := range doc.ObjectTypeDefinitions {
		names = append(names, doc.Input.ByteSliceString(def.Name))
	}
	sort.Strings(names)

	want := []string{"A", "B", "Common", "Main", "Specific"}
	if len(names) != len(want) {
		t.Fatalf("want types %v, got %v", want, names)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("want types %v, got %v", want, names)
		}
	}
}

func TestParseGraphqlDocumentFileImportPatternWithoutMatches(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "graphql-imports")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	mainPath := filepath.Join(tempDir, "main.graphql")
	if err := ioutil.WriteFile(mainPath, []byte("#import \"./missing/*.graphql\"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tempDir, "missing"), 0755); err != nil {
		t.Fatal(err)
	}

	_, report := ParseGraphqlDocumentFile(mainPath)
	if !report.HasErrors() {
		t.Fatal("expected error for import pattern without matches")
	}
}

func TestParseGraphqlDocumentFileInvalidImportPattern(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "graphql-imports")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	mainPath := filepath.Join(tempDir, "main.graphql")
	if err := ioutil.WriteFile(mainPath, []byte("#import \"./[.graphql\"\n"), 0644); err != nil {
		t.Fatal(err)
	}

	_, report := ParseGraphqlDocumentFile(mainPath)
	if !report.HasErrors() {
		t.Fatal("expected error for invalid import pattern")
	}
}
