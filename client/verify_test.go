package client

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/drand/drand/chain"
	"github.com/drand/drand/key"
	"github.com/drand/kyber/util/random"
)

type verifyMockClient struct {
	info    *chain.Info
	results map[uint64]*RandomData
}

func (v *verifyMockClient) Get(_ context.Context, round uint64) (Result, error) {
	r, ok := v.results[round]
	if !ok {
		return nil, errors.New("no result for round")
	}
	return r, nil
}

func (v *verifyMockClient) Watch(context.Context) <-chan Result { return nil }

func (v *verifyMockClient) Info(context.Context) (*chain.Info, error) { return v.info, nil }

func (v *verifyMockClient) RoundAt(time.Time) uint64 { return 0 }

func (v *verifyMockClient) Close() error { return nil }

func newVerifyMockClient(t *testing.T, rounds uint64) *verifyMockClient {
	t.Helper()
	secret := key.KeyGroup.Scalar().Pick(random.New())
	public := key.KeyGroup.Point().Mul(secret, nil)
	info := &chain.Info{
		PublicKey:   public,
		Period:      time.Minute,
		GenesisTime: time.Now().Unix() - int64(rounds-1)*int64(time.Minute/time.Second),
	}
	results := make(map[uint64]*RandomData)
	var prevSig []byte
	for round := uint64(1); round <= rounds; round++ {
		sig, err := key.AuthScheme.Sign(secret, chain.Message(round, prevSig))
		if err != nil {
			t.Fatal(err)
		}
		results[round] = &RandomData{
			Rnd:               round,
			Random:            chain.RandomnessFromSignature(sig),
			Sig:               sig,
			PreviousSig:      prevSig,
		}
		prevSig = sig
	}
	return &verifyMockClient{info: info, results: results}
}

func TestVerifyChain(t *testing.T) {
	c := newVerifyMockClient(t, 5)
	if err := VerifyChain(context.Background(), c); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyChainBrokenLink(t *testing.T) {
	c := newVerifyMockClient(t, 5)
	// corrupt the chain: round 3 points at the wrong previous signature.
	c.results[3].PreviousSig = c.results[1].Sig
	if err := VerifyChain(context.Background(), c); err == nil {
		t.Fatal("expected chain verification to fail")
	}
}

func TestVerifyChainBadSignature(t *testing.T) {
	c := newVerifyMockClient(t, 5)
	// replace the key the chain is verified against.
	other := key.KeyGroup.Scalar().Pick(random.New())
	c.info.PublicKey = key.KeyGroup.Point().Mul(other, nil)
	if err := VerifyChain(context.Background(), c); err == nil {
		t.Fatal("expected chain verification to fail")
	}
}

func TestVerifyChainMissingPreviousSignature(t *testing.T) {
	c := newVerifyMockClient(t, 5)
	c.results[2] = &RandomData{Rnd: 2, Sig: c.results[2].Sig}
	if err := VerifyChain(context.Background(), c); err == nil {
		t.Fatal("expected chain verification to fail")
	}
}
