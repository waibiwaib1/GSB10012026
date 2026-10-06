package client

import (
	"bytes"
	"context"
	"fmt"

	"github.com/drand/drand/chain"
)

// VerifyChain downloads and verifies every randomness beacon from the latest
// round back to round 1. Each beacon must verify against the chain public key,
// and its signature must be the previous signature used by the next beacon.
func VerifyChain(ctx context.Context, c Client) error {
	if c == nil {
		return fmt.Errorf("client is nil")
	}

	info, err := c.Info(ctx)
	if err != nil {
		return fmt.Errorf("getting chain info: %w", err)
	}
	if info == nil || info.PublicKey == nil {
		return fmt.Errorf("chain info is missing public key")
	}

	latest, err := c.Get(ctx, 0)
	if err != nil {
		return fmt.Errorf("getting latest beacon: %w", err)
	}
	if latest == nil {
		return fmt.Errorf("latest beacon is nil")
	}
	if latest.Round() == 0 {
		return fmt.Errorf("latest beacon has round 0")
	}

	var expectedPreviousSignature []byte
	for expectedRound := latest.Round(); expectedRound > 0; expectedRound-- {
		var result Result
		if expectedRound == latest.Round() {
			result = latest
		} else {
			result, err = c.Get(ctx, expectedRound)
			if err != nil {
				return fmt.Errorf("getting beacon at round %d: %w", expectedRound, err)
			}
		}

		if result == nil {
			return fmt.Errorf("beacon at round %d is nil", expectedRound)
		}
		if result.Round() != expectedRound {
			return fmt.Errorf("expected beacon at round %d, got round %d", expectedRound, result.Round())
		}

		previousSignature := result.PreviousSignature()
		switch {
		case expectedRound == 1 && len(previousSignature) != 0:
			return fmt.Errorf("beacon at round 1 must not have a previous signature")
		case expectedRound > 1 && len(previousSignature) == 0:
			return fmt.Errorf("beacon at round %d is missing previous signature", expectedRound)
		}

		if err := chain.Verify(info.PublicKey, previousSignature, result.Signature(), expectedRound); err != nil {
			return fmt.Errorf("verifying beacon at round %d: %w", expectedRound, err)
		}
		if !bytes.Equal(result.Randomness(), chain.RandomnessFromSignature(result.Signature())) {
			return fmt.Errorf("randomness at round %d does not match signature", expectedRound)
		}

		if expectedRound < latest.Round() && !bytes.Equal(result.Signature(), expectedPreviousSignature) {
			return fmt.Errorf("signature at round %d is not the previous signature of round %d", expectedRound, expectedRound+1)
		}

		expectedPreviousSignature = previousSignature
	}

	return nil
}
