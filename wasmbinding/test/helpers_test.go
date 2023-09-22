package wasmbinding

import (
	"testing"
	"time"

	"github.com/quicksilver-zone/quicksilver/app"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CreateTestInput(t *testing.T) (*app.Quicksilver, sdk.Context) {
	t.Helper()

	quicksilverApp := app.Setup(t, false)
	ctx := quicksilverApp.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "quicksilver-1", Time: time.Now().UTC()})
	return quicksilverApp, ctx
}

// we need to make this deterministic (same every test run), as content might affect gas costs.
func keyPubAddr() (key crypto.PrivKey, pub crypto.PubKey, addr sdk.AccAddress) { //nolint:unparam
	key = ed25519.GenPrivKey()
	pub = key.PubKey()
	addr = sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func RandomAccountAddress() sdk.AccAddress {
	_, _, addr := keyPubAddr()
	return addr
}

func RandomBech32AccountAddress() string {
	return RandomAccountAddress().String()
}
