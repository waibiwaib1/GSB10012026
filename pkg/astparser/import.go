package astparser

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/jensneuse/graphql-go-tools/pkg/ast"
	"github.com/jensneuse/graphql-go-tools/pkg/operationreport"
)

func ParseGraphqlDocumentFile(filePath string) (ast.Document, operationreport.Report) {
	input, err := mergeGraphQLFiles(filePath)
	if err != nil {
		report := operationreport.Report{}
		report.AddInternalError(err)
		return ast.Document{}, report
	}
	return ParseGraphqlDocumentBytes(input)
}

func mergeGraphQLFiles(filePath string) ([]byte, error) {
	rootPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	visited := make(map[string]bool)
	merged := make([]byte, 0)

	var merge func(string) error
	merge = func(path string) error {
		if visited[path] {
			return nil
		}

		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		visited[path] = true
		if len(merged) > 0 {
			merged = append(merged, '\n')
		}
		merged = append(merged, contents...)

		imports, err := parseImportPaths(contents)
		if err != nil {
			return fmt.Errorf("%s: %v", path, err)
		}

		directory := filepath.Dir(path)
		for _, importPath := range imports {
			matches, err := filepath.Glob(filepath.Join(directory, importPath))
			if err != nil {
				return fmt.Errorf("%s: invalid import pattern %q: %v", path, importPath, err)
			}
			if len(matches) == 0 && !hasGlobMeta(importPath) {
				return fmt.Errorf("%s: import %q: no such file or directory", path, importPath)
			}

			for _, match := range matches {
				matchPath, err := filepath.Abs(match)
				if err != nil {
					return err
				}
				if err := merge(matchPath); err != nil {
					return err
				}
			}
		}

		return nil
	}

	return merged, merge(rootPath)
}

func parseImportPaths(contents []byte) ([]string, error) {
	scanner := bufio.NewScanner(bytes.NewReader(contents))
	paths := make([]string, 0)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "#") {
			break
		}

		fields := strings.Fields(line)
		if len(fields) < 2 || fields[0] != "#import" {
			continue
		}
		if len(fields) != 2 {
			return nil, fmt.Errorf("invalid import statement %q", line)
		}

		path := strings.Trim(fields[1], `"'`)
		if path == "" {
			return nil, fmt.Errorf("invalid import statement %q", line)
		}
		paths = append(paths, path)
	}

	return paths, scanner.Err()
}

func hasGlobMeta(path string) bool {
	return strings.ContainsAny(path, `*?[`)
}
