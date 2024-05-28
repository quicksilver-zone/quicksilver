package types

import (
	"context"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/multierror"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	interchainstaking "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

func GetMappedAddresses(ctx context.Context, address string, connections []prewards.ConnectionProtocolData, config *Config) (map[string]string, error) {
	host := config.Chains[config.SourceChain]
	client, err := NewRPCClient(host, 30*time.Second)
	if err != nil {
		return nil, err
	}

	maRequest := &interchainstaking.QueryMappedAccountsRequest{
		Address: address,
	}

	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	interchainstaking.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	bytes := marshaler.MustMarshal(maRequest)
	abciquery, err := client.ABCIQuery(
		ctx,
		"/quicksilver.interchainstaking.v1.Query/MappedAccounts",
		bytes,
	)
	if err != nil {
		return nil, err
	}

	maResponse := interchainstaking.QueryMappedAccountsResponse{}
	err = marshaler.Unmarshal(abciquery.Response.Value, &maResponse)
	if err != nil {
		return nil, err
	}

	errs := map[string]error{}

	addressMap := map[string]string{}
	for chain, addrBytes := range maResponse.RemoteAddressMap {
		for _, connection := range connections {
			if connection.ChainID == chain {
				addressMap[chain], err = addressutils.EncodeAddressToBech32(connection.Prefix, sdk.AccAddress(addrBytes))
				if err != nil {
					errs[fmt.Sprintf("addressMap:%s", chain)] = err
				}
			}
		}
	}

	if len(errs) > 0 {
		return addressMap, multierror.New(errs)
	}
	return addressMap, nil
}

func GetHeights(connections []prewards.ConnectionProtocolData) map[string]int64 {
	out := make(map[string]int64, len(connections))
	for _, con := range connections {
		out[con.ChainID] = con.LastEpoch
	}
	return out
}

func GetZeroHeights(connections []prewards.ConnectionProtocolData) map[string]int64 {
	out := make(map[string]int64, len(connections))
	for _, con := range connections {
		out[con.ChainID] = 0
	}
	return out
}
