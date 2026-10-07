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
	"context"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/database"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/peer"
)

func TestStartMetrics(t *testing.T) {
	server := StartMetrics(
		"0.0.0.0:9999",
		&mockLogger{},
		func() float64 { return 0 },
		NewDatabaseList())
	defer server.Close()

	assert.IsType(t, &http.Server{}, server)

}

type metricsDBMock struct {
	database.DB

	options *database.DbOptions
	txID    uint64
}

func (db metricsDBMock) CurrentState() (*schema.ImmutableState, error) {
	return &schema.ImmutableState{TxId: db.txID}, nil
}

func (db metricsDBMock) GetOptions() *database.DbOptions {
	return db.options
}

func TestDBMetricsCollectorCollectsAllDatabases(t *testing.T) {
	dataDir := t.TempDir()
	dbList := NewDatabaseList()

	for _, dbName := range []string{"defaultdb", "seconddb"} {
		dbPath := filepath.Join(dataDir, dbName)
		require.NoError(t, os.MkdirAll(dbPath, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(dbPath, "data"), []byte(dbName), 0644))
		dbList.Append(metricsDBMock{
			options: database.DefaultOption().WithDbName(dbName).WithDbRootPath(dataDir),
			txID:    7,
		})
	}

	collector := newDBMetricsCollector(dbList)
	ch := make(chan prometheus.Metric, 4)
	collector.Collect(ch)
	close(ch)

	metricsByDatabase := map[string]map[string]float64{}
	for metric := range ch {
		dtoMetric := &dto.Metric{}
		require.NoError(t, metric.Write(dtoMetric))
		require.Len(t, dtoMetric.GetLabel(), 1)
		dbName := dtoMetric.GetLabel()[0].GetValue()
		if metricsByDatabase[dbName] == nil {
			metricsByDatabase[dbName] = map[string]float64{}
		}
		if metric.Desc().String() == collector.records.String() {
			metricsByDatabase[dbName]["records"] = dtoMetric.GetCounter().GetValue()
		} else {
			metricsByDatabase[dbName]["size"] = dtoMetric.GetCounter().GetValue()
		}
	}

	require.Len(t, metricsByDatabase, 2)
	for _, dbName := range []string{"defaultdb", "seconddb"} {
		require.Equal(t, float64(7), metricsByDatabase[dbName]["records"])
		require.Equal(t, float64(len(dbName)), metricsByDatabase[dbName]["size"])
	}
}

func TestMetricsCollection_UpdateClientMetrics(t *testing.T) {
	mc := MetricsCollection{
		UptimeCounter: prometheus.NewCounterFunc(prometheus.CounterOpts{}, func() float64 {
			return 0
		}),
		RPCsPerClientCounters: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test",
			},
			[]string{"test"},
		),
		LastMessageAtPerClientGauges: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "namespace_test",
				Subsystem: "subsystem_test",
				Name:      "test",
				Help:      "test",
			},
			[]string{
				// Which user has requested the operation?
				"test",
			},
		),
	}
	ip := net.IP{}
	ip.UnmarshalText([]byte(`127.0.0.1`))
	p := &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   ip,
			Port: 9999,
			Zone: "zone",
		},
	}
	ctx := peer.NewContext(context.TODO(), p)
	mc.UpdateClientMetrics(ctx)

	assert.IsType(t, MetricsCollection{}, mc)
}
