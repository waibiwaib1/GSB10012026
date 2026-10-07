package create

import (
	"net/http"
	"testing"

	"github.com/cli/cli/v2/internal/ghrepo"
	"github.com/cli/cli/v2/pkg/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAutoLinkCreator_Create(t *testing.T) {
	tests := []struct {
		name             string
		keyPrefix        string
		urlTemplate      string
		numeric          bool
		status           int
		wantErr          bool
		wantAlphanumeric bool
	}{
		{
			name:             "create alphanumeric autolink",
			keyPrefix:        "TICKET-",
			urlTemplate:      "https://example.com/TICKET?query=<num>",
			numeric:          false,
			status:           201,
			wantAlphanumeric: true,
		},
		{
			name:             "create numeric autolink",
			keyPrefix:        "DISCORD-",
			urlTemplate:      "https://discord.com/channels/<num>",
			numeric:          true,
			status:           201,
			wantAlphanumeric: false,
		},
		{
			name:        "http error",
			keyPrefix:   "TICKET-",
			urlTemplate: "https://example.com/TICKET?query=<num>",
			status:      404,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedPayload map[string]interface{}
			reg := &httpmock.Registry{}
			responder := httpmock.RESTPayload(tt.status, `{}`, func(payload map[string]interface{}) {
				capturedPayload = payload
			})
			if tt.wantErr {
				responder = httpmock.StatusJSONResponse(tt.status, map[string]string{
					"message": "Must have admin rights to Repository.",
				})
			}
			reg.Register(
				httpmock.REST("POST", "repos/OWNER/REPO/autolinks"),
				responder,
			)
			defer reg.Verify(t)

			autolinkCreator := &AutolinkCreator{
				HTTPClient: &http.Client{Transport: reg},
			}
			err := autolinkCreator.Create(ghrepo.New("OWNER", "REPO"), tt.keyPrefix, tt.urlTemplate, tt.numeric)

			if tt.wantErr {
				require.Error(t, err)
				assert.EqualError(t, err, "error creating autolink: HTTP 404: Must have admin rights to Repository. (https://api.github.com/repos/OWNER/REPO/autolinks)")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.keyPrefix, capturedPayload["key_prefix"])
				assert.Equal(t, tt.urlTemplate, capturedPayload["url_template"])
				assert.Equal(t, tt.wantAlphanumeric, capturedPayload["is_alphanumeric"])
			}
		})
	}
}
