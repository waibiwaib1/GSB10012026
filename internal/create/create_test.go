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
	"testing"

	"github.com/onflow/flow-cli/internal/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-cli/internal/util"
)

func TestCreateScript(t *testing.T) {
	_, _, rw := util.TestMocks(t)

	result, err := createScript([]string{"ReadNFT"}, command.GlobalFlags{}, util.NoLogger, rw, nil)
	require.NoError(t, err)
	assert.Equal(t, "cadence/scripts/ReadNFT.cdc", result.Oneliner())

	content, err := rw.ReadFile("cadence/scripts/ReadNFT.cdc")
	require.NoError(t, err)
	assert.Equal(t, scriptBoilerplate, string(content))
}

func TestCreateTransaction(t *testing.T) {
	_, _, rw := util.TestMocks(t)

	result, err := createTransaction([]string{"CreateNFT.cdc"}, command.GlobalFlags{}, util.NoLogger, rw, nil)
	require.NoError(t, err)
	assert.Equal(t, "cadence/transactions/CreateNFT.cdc", result.Oneliner())

	content, err := rw.ReadFile("cadence/transactions/CreateNFT.cdc")
	require.NoError(t, err)
	assert.Equal(t, transactionBoilerplate, string(content))
}

func TestCreateExistingFile(t *testing.T) {
	_, _, rw := util.TestMocks(t)
	path := "cadence/scripts/ReadNFT.cdc"
	require.NoError(t, rw.WriteFile(path, []byte("existing"), 0644))

	_, err := createScript([]string{"ReadNFT"}, command.GlobalFlags{}, util.NoLogger, rw, nil)
	assert.EqualError(t, err, "file already exists: cadence/scripts/ReadNFT.cdc")
}

func TestCreateRejectsPathOutsideTargetDirectory(t *testing.T) {
	_, _, rw := util.TestMocks(t)

	_, err := createScript([]string{"../ReadNFT"}, command.GlobalFlags{}, util.NoLogger, rw, nil)
	assert.EqualError(t, err, "invalid file name: ../ReadNFT.cdc")
}
