package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/drand/drand/chain"
)

// previousSignatureResult is implemented by results that expose the
// signature of the previous round, as needed to verify the whole chain.
type previousSignatureResult interface {
	PreviousSignature() []byte
}

// VerifyChain verifies the whole chain of randomness produced by the given
// client, from round 1 up to the latest round available at the time of the
// call. Each beacon's signature is verified against the chain public key, and
// the previous signature of each beacon is checked to be the signature of the
// previous round.
func VerifyChain(ctx context.Context, c Client) error {
	info, err := c.Info(ctx)
	if err != nil {
		return fmt.Errorf("getting chain info: %w", err)
	}
	last := chain.CurrentRound(time.Now().Unix(), info.Period, info.GenesisTime)
	var prevSig []byte
	for round := uint64(1); round <= last; round++ {
		res, err := c.Get(ctx, round)
		if err != nil {
			return fmt.Errorf("getting round %d: %w", round, err)
		}
		if res.Round() != round {
			return fmt.Errorf("unexpected round: got %d, expected %d", res.Round(), round)
		}
		prev, err := previousSignature(res)
		if err != nil {
			return fmt.Errorf("round %d: %w", round, err)
		}
		if !bytes.Equal(prev, prevSig) {
			return fmt.Errorf("round %d: previous signature does not match the signature of round %d", round, round-1)
		}
		if err := chain.Verify(info.PublicKey, prev, res.Signature(), round); err != nil {
			return fmt.Errorf("round %d: invalid signature: %w", round, err)
		}
		prevSig = res.Signature()
	}
	return nil
}

func previousSignature(r Result) ([]byte, error) {
	if pr, ok := r.(previousSignatureResult); ok {
		return pr.PreviousSignature(), nil
	}
	return nil, errors.New("result does not expose the previous signature")
}
