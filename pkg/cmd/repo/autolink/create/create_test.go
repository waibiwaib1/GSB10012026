package create

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/cli/cli/v2/internal/ghrepo"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCmdCreate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantOpts createOptions
		wantErr  bool
		errMsg   string
	}{
		{
			name:  "alphanumeric autolink",
			input: `"TICKET-" "https://example.com/TICKET?query=<num>"`,
			wantOpts: createOptions{
				KeyPrefix:   "TICKET-",
				URLTemplate: "https://example.com/TICKET?query=<num>",
				Numeric:     false,
			},
		},
		{
			name:  "numeric autolink",
			input: `"DISCORD-" "https://discord.com/channels/<num>" --numeric`,
			wantOpts: createOptions{
				KeyPrefix:   "DISCORD-",
				URLTemplate: "https://discord.com/channels/<num>",
				Numeric:     true,
			},
		},
		{
			name:    "no arguments",
			input:   "",
			wantErr: true,
			errMsg:  "accepts 2 arg(s), received 0",
		},
		{
			name:    "too few arguments",
			input:   `"TICKET-"`,
			wantErr: true,
			errMsg:  "accepts 2 arg(s), received 1",
		},
		{
			name:    "too many arguments",
			input:   `"TICKET-" "https://example.com/<num>" extra`,
			wantErr: true,
			errMsg:  "accepts 2 arg(s), received 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ios, _, _, _ := iostreams.Test()
			f := &cmdutil.Factory{
				IOStreams: ios,
			}
			f.HttpClient = func() (*http.Client, error) {
				return &http.Client{}, nil
			}

			argv, err := shlex.Split(tt.input)
			require.NoError(t, err)

			var gotOpts *createOptions
			cmd := NewCmdCreate(f, func(opts *createOptions) error {
				gotOpts = opts
				return nil
			})

			cmd.SetArgs(argv)
			cmd.SetIn(&bytes.Buffer{})
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})

			_, err = cmd.ExecuteC()
			if tt.wantErr {
				require.EqualError(t, err, tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantOpts.KeyPrefix, gotOpts.KeyPrefix)
				assert.Equal(t, tt.wantOpts.URLTemplate, gotOpts.URLTemplate)
				assert.Equal(t, tt.wantOpts.Numeric, gotOpts.Numeric)
			}
		})
	}
}

type stubAutoLinkCreator struct {
	err error
}

func (g stubAutoLinkCreator) Create(repo ghrepo.Interface, keyPrefix, urlTemplate string, numeric bool) error {
	return g.err
}

type testAutolinkClientCreateError struct{}

func (e testAutolinkClientCreateError) Error() string {
	return "autolink client create error"
}

func TestCreateRun(t *testing.T) {
	tests := []struct {
		name        string
		opts        *createOptions
		isTTY       bool
		stubCreator stubAutoLinkCreator
		expectedErr error
		wantStdout  string
		wantStderr  string
	}{
		{
			name: "create",
			opts: &createOptions{
				KeyPrefix:   "TICKET-",
				URLTemplate: "https://example.com/TICKET?query=<num>",
			},
			isTTY:       true,
			stubCreator: stubAutoLinkCreator{},
			wantStdout:  "✓ Autolink \"TICKET-\" created in OWNER/REPO\n",
		},
		{
			name: "create non-tty",
			opts: &createOptions{
				KeyPrefix:   "TICKET-",
				URLTemplate: "https://example.com/TICKET?query=<num>",
			},
			isTTY:       false,
			stubCreator: stubAutoLinkCreator{},
			wantStdout:  "",
		},
		{
			name: "missing num placeholder",
			opts: &createOptions{
				KeyPrefix:   "TICKET-",
				URLTemplate: "https://example.com/TICKET?query=1",
			},
			isTTY:       true,
			stubCreator: stubAutoLinkCreator{},
			expectedErr: cmdutil.FlagErrorf("URL template must contain <num>"),
		},
		{
			name: "client error",
			opts: &createOptions{
				KeyPrefix:   "TICKET-",
				URLTemplate: "https://example.com/TICKET?query=<num>",
			},
			isTTY: true,
			stubCreator: stubAutoLinkCreator{
				err: testAutolinkClientCreateError{},
			},
			expectedErr: testAutolinkClientCreateError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ios, _, stdout, _ := iostreams.Test()
			ios.SetStdoutTTY(tt.isTTY)
			ios.SetStdinTTY(tt.isTTY)
			ios.SetStderrTTY(tt.isTTY)

			opts := tt.opts
			opts.IO = ios
			opts.BaseRepo = func() (ghrepo.Interface, error) { return ghrepo.New("OWNER", "REPO"), nil }
			opts.AutolinkClient = &tt.stubCreator

			err := createRun(opts)

			if tt.expectedErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantStdout, stdout.String())
			}
		})
	}
}
