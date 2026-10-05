/*
Copyright (c) 2026 VMware, Inc. All Rights Reserved.

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

package vcenter_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vapi/vcenter"
	"github.com/vmware/govmomi/vim25/soap"
)

func TestMountISOImage(t *testing.T) {
	var request vcenter.ISOMountSpec

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/api/vcenter/iso/image", r.URL.Path)
		require.Equal(t, "mount", r.URL.Query().Get("action"))
		require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		_, err := w.Write([]byte(`"16000"`))
		require.NoError(t, err)
	}))
	defer server.Close()

	url, err := soap.ParseURL(server.URL)
	require.NoError(t, err)

	client := &rest.Client{Client: soap.NewClient(url, true)}
	cdrom, err := vcenter.NewManager(client).MountISOImage(context.Background(), vcenter.ISOMountSpec{
		LibraryItem: "item-1",
		VM:          "vm-1",
	})
	require.NoError(t, err)
	require.Equal(t, "16000", cdrom)
	require.Equal(t, "item-1", request.LibraryItem)
	require.Equal(t, "vm-1", request.VM)
}

func TestUnmountISOImage(t *testing.T) {
	var request vcenter.ISOUnmountSpec

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/api/vcenter/iso/image", r.URL.Path)
		require.Equal(t, "unmount", r.URL.Query().Get("action"))
		require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	url, err := soap.ParseURL(server.URL)
	require.NoError(t, err)

	client := &rest.Client{Client: soap.NewClient(url, true)}
	err = vcenter.NewManager(client).UnmountISOImage(context.Background(), vcenter.ISOUnmountSpec{
		VM:    "vm-1",
		CDROM: "16000",
	})
	require.NoError(t, err)
	require.Equal(t, "vm-1", request.VM)
	require.Equal(t, "16000", request.CDROM)
}
