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
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/peer"
)

func TestStartMetrics(t *testing.T) {
	server := StartMetrics(
		"0.0.0.0:9999",
		&mockLogger{},
		func() float64 { return 0 },
		&ImmuServer{})
	defer server.Close()

	assert.IsType(t, &http.Server{}, server)

}

func TestDBMetricsCollector(t *testing.T) {
	collector := newDBMetricsCollector(func() []DBMetrics {
		return []DBMetrics{{
			Name:     "defaultdb",
			Size:     42,
			NEntries: 7,
		}}
	})

	ch := make(chan prometheus.Metric, 2)
	collector.Collect(ch)
	assert.Len(t, ch, 2)
}

func TestScrapePerDBMetrics(t *testing.T) {
	collector := newDBMetricsCollector(func() []DBMetrics {
		return []DBMetrics{
			{Name: "defaultdb", Size: 10, NEntries: 2},
			{Name: "mydb", Size: 20, NEntries: 3},
		}
	})

	expected := `
# HELP immudb_db_size_bytes Database size in bytes.
# TYPE immudb_db_size_bytes gauge
immudb_db_size_bytes{db="defaultdb"} 10
immudb_db_size_bytes{db="mydb"} 20
# HELP immudb_number_of_stored_entries Number of transactions stored in each database.
# TYPE immudb_number_of_stored_entries counter
immudb_number_of_stored_entries{db="defaultdb"} 2
immudb_number_of_stored_entries{db="mydb"} 3
`
	assert.NoError(t, testutil.CollectAndCompare(collector, strings.NewReader(expected),
		"immudb_db_size_bytes", "immudb_number_of_stored_entries"))
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
