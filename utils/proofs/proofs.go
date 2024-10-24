package proofs

import (
	"fmt"

	"github.com/cosmos/gogoproto/proto"

	squareshare "github.com/celestiaorg/go-square/v2/share"

	celestiatypes "github.com/quicksilver-zone/quicksilver/third-party-chains/celestia-types/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type InclusionProof interface {
	proto.Message

	Validate(dataHash []byte) ([]byte, error)
}

var _ InclusionProof = &TendermintProof{}
var _ InclusionProof = &CelestiaProof{}

func (p *CelestiaProof) Validate(dataHash []byte) ([]byte, error) {

	shareProof, err := celestiatypes.ShareProofFromProto(*p.ShareProof)
	if err != nil {
		return nil, fmt.Errorf("unable to convert shareProof from proto: %w", err)
	}

	shares := []squareshare.Share{}
	for i, share := range shareProof.Data {
		sh, err := squareshare.NewShare(share)
		if err != nil {
			return nil, fmt.Errorf("unable to parse share %d: %w", i, err)
		}
		shares = append(shares, *sh)
	}

	txs, err := squareshare.ParseTxs(shares)
	if err != nil {
		return nil, fmt.Errorf("unable to parse txs from shareProof: %w", err)
	}

	if !shareProof.VerifyProof() {
		return nil, fmt.Errorf("share proof failed to verify")
	}

	if err := shareProof.Validate(dataHash); err != nil {
		return nil, fmt.Errorf("unable to validate celestia share proof: %w", err)
	}

	return txs[p.Index], nil
}

func (p *TendermintProof) Validate(dataHash []byte) ([]byte, error) {
	tmproof, err := tmtypes.TxProofFromProto(*p.TxProof)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal proof: %w", err)
	}
	err = tmproof.Validate(dataHash)
	if err != nil {
		return nil, fmt.Errorf("unable to validate proof: %w", err)
	}

	return tmproof.Data, nil
}
