package client_test

import (
	"context"
	"testing"
	"time"

	"github.com/drand/drand/chain"
	"github.com/drand/drand/client"
	mockresult "github.com/drand/drand/client/test/result/mock"
	"github.com/drand/drand/key"
	"github.com/drand/drand/test"
	"github.com/drand/kyber"
)

type verifyChainClient struct {
	info    *chain.Info
	results map[uint64]*mockresult.Result
	secret  kyber.Scalar
}

func (c *verifyChainClient) Get(_ context.Context, round uint64) (client.Result, error) {
	if round == 0 {
		round = c.latestRound()
	}
	result, ok := c.results[round]
	if !ok {
		return nil, nil
	}
	return result, nil
}

func (c *verifyChainClient) Watch(context.Context) <-chan client.Result {
	ch := make(chan client.Result)
	close(ch)
	return ch
}

func (c *verifyChainClient) Info(context.Context) (*chain.Info, error) {
	return c.info, nil
}

func (c *verifyChainClient) RoundAt(time.Time) uint64 {
	return c.latestRound()
}

func (c *verifyChainClient) Close() error {
	return nil
}

func (c *verifyChainClient) latestRound() uint64 {
	var latest uint64
	for round := range c.results {
		if round > latest {
			latest = round
		}
	}
	return latest
}

func TestVerifyChain(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*verifyChainClient)
		wantErr bool
	}{
		{name: "valid"},
		{
			name: "invalid signature",
			mutate: func(c *verifyChainClient) {
				c.results[2].Sig = append([]byte(nil), c.results[2].Sig...)
				c.results[2].Sig[0] ^= 1
			},
			wantErr: true,
		},
		{
			name: "missing previous signature",
			mutate: func(c *verifyChainClient) {
				c.results[2].PrevSig = nil
			},
			wantErr: true,
		},
		{
			name: "wrong randomness",
			mutate: func(c *verifyChainClient) {
				c.results[2].Rand = []byte("invalid randomness")
			},
			wantErr: true,
		},
		{
			name: "broken signature link",
			mutate: func(c *verifyChainClient) {
				previous := []byte("unlinked previous signature")
				c.results[3].PrevSig = previous
				c.results[3].Sig = signBeacon(t, c.secret, 3, previous)
				c.results[3].Rand = chain.RandomnessFromSignature(c.results[3].Sig)
			},
			wantErr: true,
		},
		{
			name: "missing first round",
			mutate: func(c *verifyChainClient) {
				delete(c.results, 1)
			},
			wantErr: true,
		},
		{
			name: "first round has previous signature",
			mutate: func(c *verifyChainClient) {
				previous := []byte("unexpected genesis previous signature")
				c.results[1].PrevSig = previous
				c.results[1].Sig = signBeacon(t, c.secret, 1, previous)
				c.results[1].Rand = chain.RandomnessFromSignature(c.results[1].Sig)
			},
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			c := signedChainClient(t, 3)
			if testCase.mutate != nil {
				testCase.mutate(c)
			}

			err := client.VerifyChain(context.Background(), c)
			if testCase.wantErr && err == nil {
				t.Fatal("expected chain verification to fail")
			}
			if !testCase.wantErr && err != nil {
				t.Fatal("expected chain verification to succeed", err)
			}
		})
	}
}

func signedChainClient(t *testing.T, latest uint64) *verifyChainClient {
	t.Helper()
	pair := test.GenerateIDs(1)[0]
	c := &verifyChainClient{
		secret: pair.Key,
		info: &chain.Info{
			PublicKey:   pair.Public.Key,
			Period:      time.Second,
			GenesisTime: time.Now().Unix(),
		},
		results: make(map[uint64]*mockresult.Result, latest),
	}

	var previous []byte
	for round := uint64(1); round <= latest; round++ {
		signature := signBeacon(t, c.secret, round, previous)
		c.results[round] = &mockresult.Result{
			Rnd:     round,
			Sig:     signature,
			PrevSig: previous,
			Rand:    chain.RandomnessFromSignature(signature),
		}
		previous = signature
	}

	return c
}

func signBeacon(t *testing.T, secret kyber.Scalar, round uint64, previous []byte) []byte {
	t.Helper()
	signature, err := key.AuthScheme.Sign(secret, chain.Message(round, previous))
	if err != nil {
		t.Fatal(err)
	}
	return signature
}
