package proofs

import (
	"encoding/hex"
	"fmt"
	"strings"

	squareshare "github.com/celestiaorg/go-square/v2/share"
	"github.com/cosmos/gogoproto/proto"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec/types"

	celestiatypes "github.com/quicksilver-zone/quicksilver/third-party-chains/celestia-types/types"
)

type InclusionProof interface {
	proto.Message

	Validate(dataHash []byte, txHash string) ([]byte, error)
}

var (
	_ InclusionProof = &TendermintProof{}
	_ InclusionProof = &CelestiaProof{}
)

func (p *CelestiaProof) Validate(dataHash []byte, txHash string) ([]byte, error) {
	if p.ShareProof == nil {
		return nil, fmt.Errorf("ShareProof is nil")
	}
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

	for _, tx := range txs {
		hash := tmhash.Sum(tx)
		hashStr := hex.EncodeToString(hash)
		if strings.EqualFold(hashStr, txHash) {
			return tx, nil
		}
	}

	return nil, fmt.Errorf("unable to find tx with hash: %s", txHash)
}

func (p *TendermintProof) Validate(dataHash []byte, txHash string) ([]byte, error) {
    if p.TxProof == nil {
        return nil, fmt.Errorf("TxProof is nil")
    }
    tmproof, err := tmtypes.TxProofFromProto(*p.TxProof)
    if err != nil {
        return nil, fmt.Errorf("unable to marshal proof: %w", err)
    }
	err = tmproof.Validate(dataHash)
	if err != nil {
		return nil, fmt.Errorf("unable to validate proof: %w", err)
	}

	hash := tmhash.Sum(tmproof.Data)
	hashStr := hex.EncodeToString(hash)
	if strings.EqualFold(hashStr, txHash) {
		return tmproof.Data, nil
	}

	return nil, fmt.Errorf("unable to find tx with hash: %s", txHash)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*InclusionProof)(nil),
		&TendermintProof{},
		&CelestiaProof{},
	)
}
