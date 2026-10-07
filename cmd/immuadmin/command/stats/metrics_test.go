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

package stats

import (
	"bytes"
	"testing"

	"github.com/codenotary/immudb/cmd/immuadmin/command/stats/statstest"
	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/require"
)

func TestMetricsPopulatesDefaultDatabaseFromLabeledMetrics(t *testing.T) {
	parser := expfmt.TextParser{}
	families, err := parser.TextToMetricFamilies(bytes.NewReader(statstest.StatsResponse))
	require.NoError(t, err)

	ms := &metrics{}
	ms.populateFrom(&families)

	require.Equal(t, "data/defaultdb", ms.db.name)
	require.Equal(t, uint64(4096), ms.db.totalBytes)
	require.Equal(t, uint64(2), ms.db.nbEntries)
}
