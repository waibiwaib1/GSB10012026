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
	"strings"
	"sync"

	"google.golang.org/grpc/peer"

	"github.com/codenotary/immudb/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// DBMetricsProvider provides metrics collected for each database.
type DBMetricsProvider interface {
	DBMetrics() []DBMetrics
}

// DBMetrics holds the metrics related to a single database.
type DBMetrics struct {
	Name     string
	Size     float64
	NEntries float64
}

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

var (
	metricsSetUptimeOnce sync.Once
	metricsSetDBsOnce    sync.Once
)

// dbMetricsCollector exposes per-database metrics. Since the set of databases
// is only known at scrape time (databases can be created/removed at runtime),
// a custom collector is used instead of a fixed set of counters/gauges.
type dbMetricsCollector struct {
	dbSizeDesc    *prometheus.Desc
	nEntriesDesc  *prometheus.Desc
	metricsGetter func() []DBMetrics
}

func newDBMetricsCollector(metricsGetter func() []DBMetrics) prometheus.Collector {
	return &dbMetricsCollector{
		dbSizeDesc: prometheus.NewDesc(
			prometheus.BuildFQName(metricsNamespace, "", "db_size_bytes"),
			"Database size in bytes.",
			[]string{"db"},
			nil,
		),
		nEntriesDesc: prometheus.NewDesc(
			prometheus.BuildFQName(metricsNamespace, "", "number_of_stored_entries"),
			"Number of transactions stored in each database.",
			[]string{"db"},
			nil,
		),
		metricsGetter: metricsGetter,
	}
}

// Describe implements prometheus.Collector.
func (c *dbMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.dbSizeDesc
	ch <- c.nEntriesDesc
}

// Collect implements prometheus.Collector.
func (c *dbMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	for _, dbm := range c.metricsGetter() {
		ch <- prometheus.MustNewConstMetric(c.dbSizeDesc, prometheus.GaugeValue, dbm.Size, dbm.Name)
		ch <- prometheus.MustNewConstMetric(c.nEntriesDesc, prometheus.CounterValue, dbm.NEntries, dbm.Name)
	}
}

func setUpUptimeCounter(uptimeCounter func() float64) {
	metricsSetUptimeOnce.Do(func() {
		Metrics.WithUptimeCounter(uptimeCounter)
	})
}

func setUpDBMetricsCollector(provider DBMetricsProvider) {
	metricsSetDBsOnce.Do(func() {
		prometheus.Register(newDBMetricsCollector(provider.DBMetrics))
	})
}

// StartMetrics listens and servers the HTTP metrics server in a new goroutine.
// The server is then returned and can be stopped using Close().
func StartMetrics(
	addr string,
	l logger.Logger,
	uptimeCounter func() float64,
	dbMetricsProvider DBMetricsProvider,
) *http.Server {
	setUpUptimeCounter(uptimeCounter)
	setUpDBMetricsCollector(dbMetricsProvider)
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
