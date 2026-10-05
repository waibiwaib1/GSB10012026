package astparser

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseGraphqlDocumentFileNestedImports(t *testing.T) {
	rootDir, err := ioutil.TempDir("", "graphql-imports")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(rootDir)

	writeFile(t, filepath.Join(rootDir, "main.graphql"), `#import "./common/schema.graphql"
type Root { id: ID }`)
	writeFile(t, filepath.Join(rootDir, "common", "schema.graphql"), `#import "./specific/a.graphql"
#import "./specific/*.graphql"
type Common { id: ID }`)
	writeFile(t, filepath.Join(rootDir, "common", "specific", "a.graphql"), `type SpecificA { id: ID }`)
	writeFile(t, filepath.Join(rootDir, "common", "specific", "b.graphql"), `type SpecificB { id: ID }`)

	merged, err := mergeGraphQLFiles(filepath.Join(rootDir, "main.graphql"))
	if err != nil {
		t.Fatal(err)
	}

	for _, typeName := range []string{"Root", "Common", "SpecificA", "SpecificB"} {
		if !strings.Contains(string(merged), "type "+typeName) {
			t.Fatalf("merged schema does not contain %s:\n%s", typeName, merged)
		}
	}

	doc, report := ParseGraphqlDocumentFile(filepath.Join(rootDir, "main.graphql"))
	if report.HasErrors() {
		t.Fatal(report.Error())
	}
	if got := len(doc.ObjectTypeDefinitions); got != 4 {
		t.Fatalf("expected 4 object type definitions, got %d", got)
	}
}

func TestParseGraphqlDocumentFileMissingImport(t *testing.T) {
	rootDir, err := ioutil.TempDir("", "graphql-imports")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(rootDir)

	writeFile(t, filepath.Join(rootDir, "main.graphql"), `#import "./missing.graphql"
type Root { id: ID }`)

	_, report := ParseGraphqlDocumentFile(filepath.Join(rootDir, "main.graphql"))
	if !report.HasErrors() {
		t.Fatal("expected missing import error")
	}
}

func writeFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(path, []byte(contents), 0644); err != nil {
		t.Fatal(err)
	}
}
