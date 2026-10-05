package astparser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jensneuse/graphql-go-tools/pkg/ast"
	"github.com/jensneuse/graphql-go-tools/pkg/operationreport"
)

func ParseGraphqlDocumentFile(filePath string) (ast.Document, operationreport.Report) {
	doc := *ast.NewDocument()
	report := operationreport.Report{}

	input, err := readGraphqlDocument(filePath, map[string]struct{}{})
	if err != nil {
		report.AddInternalError(err)
		return doc, report
	}

	doc.Input.ResetInputBytes(input)
	NewParser().Parse(&doc, &report)
	return doc, report
}

func readGraphqlDocument(filePath string, imported map[string]struct{}) ([]byte, error) {
	absolutePath, err := canonicalImportPath(filePath)
	if err != nil {
		return nil, err
	}

	if _, ok := imported[absolutePath]; ok {
		return nil, nil
	}
	imported[absolutePath] = struct{}{}

	input, err := ioutil.ReadFile(absolutePath)
	if err != nil {
		return nil, err
	}

	patterns, err := importPatterns(input)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", absolutePath, err)
	}

	var output bytes.Buffer
	fileDirectory := filepath.Dir(absolutePath)
	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(fileDirectory, pattern))
		if err != nil {
			return nil, err
		}
		if len(matches) == 0 {
			return nil, fmt.Errorf("import pattern matches no files: %s", pattern)
		}

		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil {
				return nil, err
			}
			if info.IsDir() {
				continue
			}

			importedInput, err := readGraphqlDocument(match, imported)
			if err != nil {
				return nil, err
			}
			if len(importedInput) > 0 {
				output.Write(importedInput)
				if importedInput[len(importedInput)-1] != '\n' {
					output.WriteByte('\n')
				}
			}
		}
	}

	output.Write(input)
	return output.Bytes(), nil
}

func canonicalImportPath(filePath string) (string, error) {
	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}
	evaluatedPath, err := filepath.EvalSymlinks(absolutePath)
	if err != nil {
		return "", err
	}
	return filepath.Clean(evaluatedPath), nil
}

func importPatterns(input []byte) ([]string, error) {
	var patterns []string

	for _, rawLine := range bytes.Split(input, []byte{'\n'}) {
		line := bytes.TrimSpace(bytes.TrimSuffix(rawLine, []byte{'\r'}))
		if len(line) == 0 {
			continue
		}
		if line[0] != '#' {
			break
		}

		comment := bytes.TrimSpace(line[1:])
		if !bytes.HasPrefix(comment, []byte("import")) {
			continue
		}
		if len(comment) == len("import") || comment[len("import")] != ' ' && comment[len("import")] != '\t' {
			return nil, fmt.Errorf("invalid import statement: %s", line)
		}

		pattern := string(bytes.TrimSpace(comment[len("import"):]))
		if len(pattern) >= 2 {
			if (pattern[0] == '"' && pattern[len(pattern)-1] == '"') ||
				(pattern[0] == '\'' && pattern[len(pattern)-1] == '\'') {
				pattern = pattern[1 : len(pattern)-1]
			}
		}
		if pattern == "" {
			return nil, fmt.Errorf("invalid import statement: %s", line)
		}
		if _, err := filepath.Match(pattern, pattern); err != nil {
			return nil, err
		}

		patterns = append(patterns, pattern)
	}

	return patterns, nil
}
