// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package csp

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-safeweb/safehttp"
	"github.com/google/go-safeweb/safehttp/safehttptest"
)

func TestDisable(t *testing.T) {
	tests := []struct {
		name        string
		interceptor safehttp.Interceptor
	}{
		{
			name:        "Interceptor value",
			interceptor: Default(""),
		},
		{
			name: "Interceptor pointer",
			interceptor: func() safehttp.Interceptor {
				i := Default("")
				return &i
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			disabled, ok := Disable().Apply(tt.interceptor)
			if !ok {
				t.Fatal("Disable().Apply(interceptor) got ok: false want: true")
			}

			rr := safehttptest.NewResponseRecorder()
			req := safehttptest.NewRequest(safehttp.MethodGet, "/", nil)
			disabled.Before(rr.ResponseWriter, req)

			if diff := cmp.Diff(map[string][]string{}, map[string][]string(rr.Header())); diff != "" {
				t.Errorf("rr.Header() mismatch (-want +got):\n%s", diff)
			}
			if got := req.Context().Value(ctxKey{}); got != nil {
				t.Errorf("req.Context().Value(ctxKey{}) got: %v want: nil", got)
			}
			if got := rr.Status(); got != safehttp.StatusOK {
				t.Errorf("rr.Status() got: %v want: %v", got, safehttp.StatusOK)
			}
		})
	}
}

func TestDisableNoMatch(t *testing.T) {
	interceptor := staticHeadersInterceptor{}
	got, ok := Disable().Apply(interceptor)
	if ok {
		t.Error("Disable().Apply(interceptor) got ok: true want: false")
	}
	if got != safehttp.Interceptor(interceptor) {
		t.Error("Disable().Apply(interceptor) returned a modified interceptor")
	}
}

type staticHeadersInterceptor struct{}

func (staticHeadersInterceptor) Before(w *safehttp.ResponseWriter, r *safehttp.IncomingRequest) safehttp.Result {
	return safehttp.Result{}
}
