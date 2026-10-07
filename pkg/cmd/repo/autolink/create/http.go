package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cli/cli/v2/api"
	"github.com/cli/cli/v2/internal/ghinstance"
	"github.com/cli/cli/v2/internal/ghrepo"
)

type AutolinkCreator struct {
	HTTPClient *http.Client
}

func (a *AutolinkCreator) Create(repo ghrepo.Interface, keyPrefix, urlTemplate string, numeric bool) error {
	path := fmt.Sprintf("repos/%s/%s/autolinks", repo.RepoOwner(), repo.RepoName())
	url := ghinstance.RESTPrefix(repo.RepoHost()) + path

	payload := map[string]interface{}{
		"key_prefix":      keyPrefix,
		"url_template":    urlTemplate,
		"is_alphanumeric": !numeric,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return fmt.Errorf("error creating autolink: %w", api.HandleHTTPError(resp))
	}

	return nil
}
