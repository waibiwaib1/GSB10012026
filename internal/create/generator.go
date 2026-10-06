/*
 * Flow CLI
 *
 * Copyright 2019 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package create

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onflow/flow-cli/flowkit"
	"github.com/onflow/flow-cli/flowkit/output"
	"github.com/onflow/flow-cli/internal/command"
)

const (
	cadenceDir      = "cadence"
	scriptsDir      = "scripts"
	transactionsDir = "transactions"
	cadenceExt      = ".cdc"
)

const scriptBoilerplate = "pub fun main() {}\n"

const transactionBoilerplate = `transaction() {
	prepare() {}

	execute {}
}
`

type directoryCreator interface {
	MkdirAll(path string, perm os.FileMode) error
}

var scriptCommand = &command.Command{
	Cmd: &cobra.Command{
		Use:     "script <name>",
		Short:   "Create a Cadence script",
		Example: "flow create script ReadNFT",
		Args:    cobra.ExactArgs(1),
	},
	Flags: &struct{}{},
	Run:   createScript,
}

var transactionCommand = &command.Command{
	Cmd: &cobra.Command{
		Use:     "tx <name>",
		Aliases: []string{"transaction"},
		Short:   "Create a Cadence transaction",
		Example: "flow create tx CreateNFT",
		Args:    cobra.ExactArgs(1),
	},
	Flags: &struct{}{},
	Run:   createTransaction,
}

func createScript(
	args []string,
	_ command.GlobalFlags,
	_ output.Logger,
	readerWriter flowkit.ReaderWriter,
	_ flowkit.Services,
) (command.Result, error) {
	return createCadenceFile(args[0], scriptsDir, scriptBoilerplate, readerWriter)
}

func createTransaction(
	args []string,
	_ command.GlobalFlags,
	_ output.Logger,
	readerWriter flowkit.ReaderWriter,
	_ flowkit.Services,
) (command.Result, error) {
	return createCadenceFile(args[0], transactionsDir, transactionBoilerplate, readerWriter)
}

func createCadenceFile(
	name string,
	directory string,
	boilerplate string,
	readerWriter flowkit.ReaderWriter,
) (command.Result, error) {
	filename := filepath.Base(filepath.Clean(name))
	if filename == "." || filename == string(filepath.Separator) {
		return nil, fmt.Errorf("invalid Cadence filename: %s", name)
	}

	if !strings.HasSuffix(filename, cadenceExt) {
		filename += cadenceExt
	}

	path := filepath.Join(cadenceDir, directory, filename)
	if err := ensureFileDoesNotExist(path, readerWriter); err != nil {
		return nil, err
	}

	if creator, ok := readerWriter.(directoryCreator); ok {
		if err := creator.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, fmt.Errorf("failed creating directory: %w", err)
		}
	}

	if err := readerWriter.WriteFile(path, []byte(boilerplate), 0644); err != nil {
		return nil, fmt.Errorf("failed creating file: %w", err)
	}

	return &createResult{path: path}, nil
}

func ensureFileDoesNotExist(path string, readerWriter flowkit.ReaderWriter) error {
	_, err := readerWriter.ReadFile(path)
	if err == nil {
		return fmt.Errorf("file already exists: %s", path)
	}

	if !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("failed checking if file exists: %w", err)
	}

	return nil
}

type createResult struct {
	path string
}

func (r *createResult) JSON() any {
	return map[string]string{"path": r.path}
}

func (r *createResult) String() string {
	return r.path
}

func (r *createResult) Oneliner() string {
	return r.path
}
