/*
Copyright 2021 CodeNotary, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/database"
	"github.com/stretchr/testify/require"
)

type dbMock struct {
	database.DB

	currentStateF func() (*schema.ImmutableState, error)
	options       *database.DbOptions
}

func (dbm dbMock) CurrentState() (*schema.ImmutableState, error) {
	if dbm.currentStateF != nil {
		return dbm.currentStateF()
	}
	return &schema.ImmutableState{TxId: 99}, nil
}

func (dbm dbMock) GetOptions() *database.DbOptions {
	if dbm.options != nil {
		return dbm.options
	}
	return database.DefaultOption().WithDbName("defaultdb")
}

func TestMetricFuncDBRecordsCounter(t *testing.T) {
	s := ImmuServer{
		dbList: &databaseList{
			databases: []database.DB{dbMock{}},
		},
	}
	nbRecords := s.metricFuncDBRecordsCounter(s.dbList.GetByIndex(0))
	require.Equal(t, 99, int(nbRecords))
}

func TestMetricFuncServerUptimeCounter(t *testing.T) {
	s := ImmuServer{}
	s.metricFuncServerUptimeCounter()
}

func TestMetricFuncDBSize(t *testing.T) {
	s := ImmuServer{
		Options: &Options{
			Dir:           ".",
			defaultDbName: "TestMetricFuncServerUptimeCounter_DefaultDB",
		},
	}

	defaultDBPath := filepath.Join(s.Options.Dir, s.Options.defaultDbName)
	require.NoError(t, os.MkdirAll(defaultDBPath, 0777))
	defer os.RemoveAll(defaultDBPath)
	f, err := os.Create(filepath.Join(defaultDBPath, "some-db-file"))
	require.NoError(t, err)
	_, err = f.WriteString("some-db-content")
	require.NoError(t, err)
	require.NoError(t, f.Close())

	size := s.metricFuncDBSize(".", s.Options.defaultDbName)
	require.Greater(t, size, float64(0))
}

func TestDBMetrics(t *testing.T) {
	s := ImmuServer{
		Options: DefaultOptions(),
		dbList: &databaseList{
			databases: []database.DB{
				dbMock{options: database.DefaultOption().WithDbName("defaultdb")},
				dbMock{options: database.DefaultOption().WithDbName("mydb")},
			},
		},
	}

	dbMetrics := s.DBMetrics()
	require.Len(t, dbMetrics, 2)
	require.Equal(t, "defaultdb", dbMetrics[0].Name)
	require.Equal(t, "mydb", dbMetrics[1].Name)
	require.Equal(t, float64(99), dbMetrics[0].NEntries)
}
