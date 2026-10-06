// Copyright 2024 OpenPubkey
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package providers

import (
	"context"
	"net/http"
	"time"

	"github.com/openpubkey/openpubkey/discover"
)

// GitlabOptions is an options struct that configures how providers.GitlabOp
// operates. See providers.GetDefaultGitlabOpOptions for the recommended default
// values to use when interacting with Gitlab as the OpenIdProvider.
//
// This is the OP for Gitlab users (humans) authenticating through the
// browser-based OIDC authorization code flow. For the OP that supplies ID
// Tokens to Gitlab CI workflows (workloads) via an environment variable, see
// GitlabCiOp (providers/gitlab_ci.go).
type GitlabOptions struct {
	// ClientID is the client ID of the OIDC application. It should be the
	// expected "aud" claim in received ID tokens from the OP.
	ClientID string
	// Issuer is the OP's issuer URI for performing OIDC authorization and
	// discovery.
	Issuer string
	// Scopes is the list of scopes to send to the OP in the initial
	// authorization request.
	Scopes []string
	// RedirectURIs is the list of authorized redirect URIs that can be
	// redirected to by the OP after the user completes the authorization code
	// flow exchange. Ensure that your OIDC application is configured to accept
	// these URIs otherwise an error may occur.
	RedirectURIs []string
	// GQSign denotes if the received ID token should be upgraded to a GQ token
	// using GQ signatures.
	GQSign bool
	// OpenBrowser denotes if the client's default browser should be opened
	// automatically when performing the OIDC authorization flow. This value
	// should typically be set to true, unless performing some headless
	// automation (e.g. integration tests) where you don't want the browser to
	// open.
	OpenBrowser bool
	// HttpClient is the http.Client to use when making queries to the OP (OIDC
	// code exchange, refresh, verification of ID token, fetch of JWKS endpoint,
	// etc.). If nil, then http.DefaultClient is used.
	HttpClient *http.Client
	// IssuedAtOffset configures the offset to add when validating the "iss" and
	// "exp" claims of received ID tokens from the OP.
	IssuedAtOffset time.Duration
}

func GetDefaultGitlabOpOptions() *GitlabOptions {
	return &GitlabOptions{
		Issuer:   gitlabIssuer,
		ClientID: "8d8b7024572c7fd501f64374dec6bba37096783dfcd792b3988104be08cb6923",
		Scopes:   []string{"openid email"},
		RedirectURIs: []string{
			"http://localhost:3000/login-callback",
			"http://localhost:10001/login-callback",
			"http://localhost:11110/login-callback",
		},
		GQSign:         false,
		OpenBrowser:    true,
		HttpClient:     nil,
		IssuedAtOffset: 1 * time.Minute,
	}
}

// NewGitlabOp creates a Gitlab OP (OpenID Provider) using the default
// configuration options. It uses the OIDC Relying Party (Client) setup by the
// OpenPubkey project to authenticate Gitlab users (humans) via the
// browser-based OIDC authorization code flow. This is not the OP for Gitlab CI
// workflows; that functionality is provided by NewGitlabCiOp.
func NewGitlabOp() BrowserOpenIdProvider {
	options := GetDefaultGitlabOpOptions()
	return NewGitlabOpWithOptions(options)
}

// NewGitlabOpWithOptions creates a Gitlab OP with configuration specified
// using an options struct. This is useful if you want to use your own OIDC
// Client or override the configuration.
func NewGitlabOpWithOptions(opts *GitlabOptions) BrowserOpenIdProvider {
	return &StandardOp{
		clientID:                  opts.ClientID,
		Scopes:                    opts.Scopes,
		RedirectURIs:              opts.RedirectURIs,
		GQSign:                    opts.GQSign,
		OpenBrowser:               opts.OpenBrowser,
		HttpClient:                opts.HttpClient,
		IssuedAtOffset:            opts.IssuedAtOffset,
		issuer:                    opts.Issuer,
		requestTokensOverrideFunc: nil,
		publicKeyFinder: discover.PublicKeyFinder{
			JwksFunc: func(ctx context.Context, issuer string) ([]byte, error) {
				return discover.GetJwksByIssuer(ctx, issuer, opts.HttpClient)
			},
		},
	}
}

type GitlabOp = StandardOp

var _ OpenIdProvider = (*GitlabOp)(nil)
var _ BrowserOpenIdProvider = (*GitlabOp)(nil)
var _ RefreshableOpenIdProvider = (*GitlabOp)(nil)
