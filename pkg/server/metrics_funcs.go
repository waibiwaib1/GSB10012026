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
	"time"

	"github.com/codenotary/immudb/pkg/database"
)

func (s *ImmuServer) metricFuncServerUptimeCounter() float64 {
	return time.Since(startedAt).Hours()
}

// DBMetrics returns the size on disk and the number of stored transactions
// for each database opened by the server.
func (s *ImmuServer) DBMetrics() []DBMetrics {
	result := make([]DBMetrics, 0, s.dbList.Length())
	for i := 0; i < s.dbList.Length(); i++ {
		db := s.dbList.GetByIndex(int64(i))
		opts := db.GetOptions()
		result = append(result, DBMetrics{
			Name:     opts.GetDbName(),
			Size:     s.metricFuncDBSize(opts.GetDbRootPath(), opts.GetDbName()),
			NEntries: s.metricFuncDBRecordsCounter(db),
		})
	}
	return result
}

func (s *ImmuServer) metricFuncDBRecordsCounter(db database.DB) float64 {
	ic, err := db.CurrentState()
	if err != nil {
		return 0
	}
	return float64(ic.GetTxId())
}

func (s *ImmuServer) metricFuncDBSize(rootPath, dbName string) float64 {
	var dbDirSizeBytes int64 = 0
	readSize := func(path string, file os.FileInfo, err error) error {
		if err == nil && !file.IsDir() {
			dbDirSizeBytes += file.Size()
		}
		return nil
	}
	filepath.Walk(filepath.Join(rootPath, dbName), readSize)
	return float64(dbDirSizeBytes)
}
