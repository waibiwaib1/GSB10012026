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
	"strings"
	"testing"

	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/openpubkey/openpubkey/providers/mocks"
	"github.com/openpubkey/openpubkey/util"
	"github.com/stretchr/testify/require"
)

func TestGitlabUserSimpleRequest(t *testing.T) {
	issuer := gitlabIssuer
	providerOverride, err := mocks.NewMockProviderBackend(issuer, 2)
	require.NoError(t, err)

	op := &GitlabOp{
		issuer:                    gitlabIssuer,
		clientID:                  "test-client-id",
		publicKeyFinder:           providerOverride.PublicKeyFinder,
		requestTokensOverrideFunc: providerOverride.RequestTokensOverrideFunc,
	}

	cic := GenCIC(t)
	expSigningKey, expKeyID, expRecord := providerOverride.RandomSigningKey()

	idTokenTemplate := mocks.IDTokenTemplate{
		CommitFunc:  mocks.AddNonceCommit,
		Issuer:      issuer,
		Nonce:       "empty",
		NoNonce:     false,
		Aud:         "test-client-id",
		KeyID:       expKeyID,
		NoKeyID:     false,
		Alg:         expRecord.Alg,
		NoAlg:       false,
		ExtraClaims: map[string]any{"extraClaim": "extraClaimValue"},
		SigningKey:  expSigningKey,
	}
	providerOverride.SetIDTokenTemplate(&idTokenTemplate)

	tokens, err := op.RequestTokens(context.Background(), cic)
	require.NoError(t, err)
	idToken := tokens.IDToken

	cicHash, err := cic.Hash()
	require.NoError(t, err)
	require.NotNil(t, cicHash)

	_, payloadB64, _, err := jws.SplitCompact(idToken)
	require.NoError(t, err)
	payload, err := util.Base64DecodeForJWT(payloadB64)
	require.NoError(t, err)
	require.Contains(t, string(payload), string(cicHash))

	require.Equal(t, "mock-refresh-token", string(tokens.RefreshToken))
	require.Equal(t, "mock-access-token", string(tokens.AccessToken))

	require.NoError(t, op.VerifyIDToken(context.Background(), idToken, cic))
}

func TestGitlabUserDefaultOptions(t *testing.T) {
	opts := GetDefaultGitlabOpOptions()
	require.Equal(t, gitlabIssuer, opts.Issuer)
	require.NotEmpty(t, opts.ClientID)
	require.Contains(t, strings.Join(opts.Scopes, " "), "openid")
	require.NotEmpty(t, opts.RedirectURIs)

	op := NewGitlabOpWithOptions(opts)
	require.Equal(t, gitlabIssuer, op.Issuer())
	require.Equal(t, opts.ClientID, op.ClientID())
}
