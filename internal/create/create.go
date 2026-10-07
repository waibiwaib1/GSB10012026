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
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onflow/flow-cli/flowkit"
	"github.com/onflow/flow-cli/flowkit/output"
	"github.com/onflow/flow-cli/internal/command"
)

const (
	scriptBoilerplate           = "pub fun main() {}\n"
	transactionBoilerplate      = "transaction() {\n\tprepare() {}\n\n\texecute {}\n}\n"
	defaultFilePermissions      = os.FileMode(0644)
	defaultDirectoryPermissions = os.FileMode(0755)
)

var Cmd = &cobra.Command{
	Use:     "create <script|tx> <name>",
	Short:   "Create a new script or transaction",
	GroupID: "project",
}

type result struct {
	path string
}

func (r *result) JSON() any {
	return map[string]string{"path": r.path}
}

func (r *result) String() string {
	return fmt.Sprintf("Created %s\n", r.path)
}

func (r *result) Oneliner() string {
	return r.path
}

func init() {
	scriptCommand.AddToParent(Cmd)
	transactionCommand.AddToParent(Cmd)
}

var scriptCommand = &command.Command{
	Cmd: &cobra.Command{
		Use:     "script <name>",
		Short:   "Create a new script",
		Example: "flow create script ReadNFT",
		Aliases: []string{"scripts"},
		Args:    cobra.ExactArgs(1),
	},
	Run: createScript,
}

var transactionCommand = &command.Command{
	Cmd: &cobra.Command{
		Use:     "tx <name>",
		Short:   "Create a new transaction",
		Example: "flow create tx CreateNFT",
		Aliases: []string{"transaction", "transactions"},
		Args:    cobra.ExactArgs(1),
	},
	Run: createTransaction,
}

func createScript(
	args []string,
	_ command.GlobalFlags,
	_ output.Logger,
	rw flowkit.ReaderWriter,
	_ flowkit.Services,
) (command.Result, error) {
	return createFile(args[0], "cadence/scripts", scriptBoilerplate, rw)
}

func createTransaction(
	args []string,
	_ command.GlobalFlags,
	_ output.Logger,
	rw flowkit.ReaderWriter,
	_ flowkit.Services,
) (command.Result, error) {
	return createFile(args[0], "cadence/transactions", transactionBoilerplate, rw)
}

func createFile(name, directory, boilerplate string, rw flowkit.ReaderWriter) (*result, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("file name cannot be empty")
	}

	name = path.Clean(name)
	if path.Ext(name) == "" {
		name += ".cdc"
	}

	filePath := path.Join(directory, name)
	if !strings.HasPrefix(filePath, directory+"/") {
		return nil, fmt.Errorf("invalid file name: %s", name)
	}

	if _, err := rw.ReadFile(filePath); err == nil {
		return nil, fmt.Errorf("file already exists: %s", filePath)
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	if err := mkdirAll(path.Dir(filePath), defaultDirectoryPermissions, rw); err != nil {
		return nil, err
	}

	if err := rw.WriteFile(filePath, []byte(boilerplate), defaultFilePermissions); err != nil {
		return nil, err
	}

	return &result{path: filePath}, nil
}

type directoryCreator interface {
	MkdirAll(path string, perm os.FileMode) error
}

func mkdirAll(path string, perm os.FileMode, rw flowkit.ReaderWriter) error {
	if creator, ok := rw.(directoryCreator); ok {
		return creator.MkdirAll(path, perm)
	}

	return os.MkdirAll(path, perm)
}
