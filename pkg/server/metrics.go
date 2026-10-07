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
	"expvar"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"google.golang.org/grpc/peer"

	"github.com/codenotary/immudb/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsCollection immudb Prometheus metrics collection
type MetricsCollection struct {
	UptimeCounter                prometheus.CounterFunc
	RPCsPerClientCounters        *prometheus.CounterVec
	LastMessageAtPerClientGauges *prometheus.GaugeVec
}

var metricsNamespace = "immudb"

// WithUptimeCounter ...
func (mc *MetricsCollection) WithUptimeCounter(f func() float64) {
	mc.UptimeCounter = promauto.NewCounterFunc(
		prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Name:      "uptime_hours",
			Help:      "Server uptime in hours.",
		},
		f,
	)
}

type databasePathOptions interface {
	GetDbName() string
	GetDbRootPath() string
}

type dbMetricsCollector struct {
	mutex  sync.RWMutex
	dbList DatabaseList

	records *prometheus.Desc
	size    *prometheus.Desc
}

func newDBMetricsCollector(dbList DatabaseList) *dbMetricsCollector {
	return &dbMetricsCollector{
		dbList: dbList,
		records: prometheus.NewDesc(
			prometheus.BuildFQName(metricsNamespace, "", "number_of_stored_entries"),
			"Number of key-value entries currently stored by the database.",
			[]string{"database"},
			nil,
		),
		size: prometheus.NewDesc(
			prometheus.BuildFQName(metricsNamespace, "", "db_size_bytes"),
			"Database size in bytes.",
			[]string{"database"},
			nil,
		),
	}
}

func (c *dbMetricsCollector) setDatabaseList(dbList DatabaseList) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.dbList = dbList
}

func (c *dbMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.records
	ch <- c.size
}

func (c *dbMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.RLock()
	dbList := c.dbList
	c.mutex.RUnlock()

	if dbList == nil {
		return
	}

	for i := 0; i < dbList.Length(); i++ {
		db := dbList.GetByIndex(int64(i))
		options := db.GetOptions()
		dbName := options.GetDbName()

		var records float64
		if state, err := db.CurrentState(); err == nil {
			records = float64(state.GetTxId())
		}
		ch <- prometheus.MustNewConstMetric(c.records, prometheus.CounterValue, records, dbName)
		ch <- prometheus.MustNewConstMetric(c.size, prometheus.CounterValue, databaseSizeBytes(options), dbName)
	}
}

func databaseSizeBytes(options databasePathOptions) float64 {
	var size int64
	_ = filepath.Walk(filepath.Join(options.GetDbRootPath(), options.GetDbName()), func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return float64(size)
}

// UpdateClientMetrics ...
func (mc *MetricsCollection) UpdateClientMetrics(ctx context.Context) {
	p, ok := peer.FromContext(ctx)
	if ok && p != nil {
		ipAndPort := strings.Split(p.Addr.String(), ":")
		if len(ipAndPort) > 0 {
			mc.RPCsPerClientCounters.WithLabelValues(ipAndPort[0]).Inc()
			mc.LastMessageAtPerClientGauges.WithLabelValues(ipAndPort[0]).SetToCurrentTime()
		}
	}
}

// Metrics immudb Prometheus metrics collection
var Metrics = MetricsCollection{
	RPCsPerClientCounters: promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Name:      "number_of_rpcs_per_client",
			Help:      "Number of handled RPCs per client.",
		},
		[]string{"ip"},
	),
	LastMessageAtPerClientGauges: promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: metricsNamespace,
			Name:      "clients_last_message_at_unix_seconds",
			Help:      "Timestamp at which clients have sent their most recent message.",
		},
		[]string{"ip"},
	),
}

var dbMetrics = newDBMetricsCollector(nil)

func init() {
	prometheus.MustRegister(dbMetrics)
}

// StartMetrics listens and servers the HTTP metrics server in a new goroutine.
// The server is then returned and can be stopped using Close().
func StartMetrics(
	addr string,
	l logger.Logger,
	uptimeCounter func() float64,
	dbList DatabaseList,
) *http.Server {
	Metrics.WithUptimeCounter(uptimeCounter)
	dbMetrics.setDatabaseList(dbList)
	// expvar package adds a handler in to the default HTTP server (which has to be started explicitly),
	// and serves up the metrics at the /debug/vars endpoint.
	// Here we're registering both expvar and promhttp handlers in our custom server.
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/debug/vars", expvar.Handler())
	server := &http.Server{Addr: addr, Handler: mux}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				l.Debugf("Metrics http server closed")
			} else {
				l.Errorf("Metrics error: %s", err)
			}

		}
	}()

	return server
}
