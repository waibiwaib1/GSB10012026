package create

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/v2/internal/ghrepo"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type createOptions struct {
	BaseRepo       func() (ghrepo.Interface, error)
	AutolinkClient AutolinkClient
	IO             *iostreams.IOStreams

	KeyPrefix   string
	URLTemplate string
	Numeric     bool
}

type AutolinkClient interface {
	Create(repo ghrepo.Interface, keyPrefix, urlTemplate string, numeric bool) error
}

func NewCmdCreate(f *cmdutil.Factory, runF func(*createOptions) error) *cobra.Command {
	opts := &createOptions{
		IO: f.IOStreams,
	}

	cmd := &cobra.Command{
		Use:   "create <keyPrefix> <urlTemplate>",
		Short: "Create an autolink reference for a GitHub repository",
		Long: heredoc.Doc(`
			Create an autolink reference for a GitHub repository.

			The URL template must contain <num> as a placeholder for the reference number.

			Autolinks can be either alphanumeric or numeric. By default, autolinks are
			alphanumeric; use the --numeric flag to create a numeric autolink.
		`),
		Example: heredoc.Doc(`
			# Create alphanumeric autolink
			gh repo autolink create "TICKET-" "https://example.com/TICKET?query=<num>"

			# Create numeric autolink
			gh repo autolink create "DISCORD-" "https://discord.com/channels/<num>" --numeric
		`),
		Aliases: []string{"new"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.BaseRepo = f.BaseRepo
			opts.KeyPrefix = args[0]
			opts.URLTemplate = args[1]

			httpClient, err := f.HttpClient()
			if err != nil {
				return err
			}
			opts.AutolinkClient = &AutolinkCreator{HTTPClient: httpClient}

			if runF != nil {
				return runF(opts)
			}

			return createRun(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Numeric, "numeric", "n", false, "Mark autolink as non-alphanumeric")

	return cmd
}

func createRun(opts *createOptions) error {
	repo, err := opts.BaseRepo()
	if err != nil {
		return err
	}

	if !strings.Contains(opts.URLTemplate, "<num>") {
		return cmdutil.FlagErrorf("URL template must contain <num>")
	}

	err = opts.AutolinkClient.Create(repo, opts.KeyPrefix, opts.URLTemplate, opts.Numeric)
	if err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Autolink %q created in %s\n", cs.SuccessIcon(), opts.KeyPrefix, ghrepo.FullName(repo))
	}

	return nil
}
