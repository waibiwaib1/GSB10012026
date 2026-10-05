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

package vcenter

import (
	"context"
	"net/http"
)

const isoImagePath = "/api/vcenter/iso/image"

// ISOMountSpec specifies a Content Library ISO to mount on a virtual machine.
type ISOMountSpec struct {
	LibraryItem string `json:"library_item"`
	VM          string `json:"vm"`
}

// ISOUnmountSpec specifies a Content Library ISO-backed CD-ROM to unmount.
type ISOUnmountSpec struct {
	VM    string `json:"vm"`
	CDROM string `json:"cdrom"`
}

// MountISOImage mounts an ISO from a Content Library and returns the ID of the created CD-ROM.
func (c *Manager) MountISOImage(ctx context.Context, spec ISOMountSpec) (string, error) {
	url := c.Resource(isoImagePath).WithParam("action", "mount")
	req := url.Request(http.MethodPost, spec)

	var cdrom string
	return cdrom, c.Do(ctx, req, &cdrom)
}

// UnmountISOImage unmounts a previously mounted Content Library ISO-backed CD-ROM.
func (c *Manager) UnmountISOImage(ctx context.Context, spec ISOUnmountSpec) error {
	url := c.Resource(isoImagePath).WithParam("action", "unmount")
	req := url.Request(http.MethodPost, spec)

	return c.Do(ctx, req, nil)
}
