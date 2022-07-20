package interchainstaking

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	distrTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"

	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v4/modules/core/05-port/types"
	host "github.com/cosmos/ibc-go/v4/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v4/modules/core/exported"
)

var _ porttypes.IBCModule = IBCModule{}

// IBCModule implements the ICS26 interface for interchain accounts controller chains
type IBCModule struct {
	keeper keeper.Keeper
}

// NewIBCModule creates a new IBCModule given the keeper
func NewIBCModule(k keeper.Keeper) IBCModule {
	return IBCModule{
		keeper: k,
	}
}

// OnChanOpenInit implements the IBCModule interface
func (im IBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) (string, error) {
	return "", im.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID))
}

// OnChanOpenTry implements the IBCModule interface
func (im IBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {
	return "", nil
}

// OnChanOpenAck implements the IBCModule interface
func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterPartyChannelId string,
	counterpartyVersion string,
) error {
	// get connection from port
	connectionId, err := im.keeper.GetConnectionForPort(ctx, portID)
	if err != nil {
		ctx.Logger().Error("unable to get connection for port " + portID)
	}

	// get chain id from connection
	chainId, err := im.keeper.GetChainID(ctx, connectionId)
	if err != nil {
		ctx.Logger().Error(
			"Unable to obtain chain for given connection and port",
			"ConnectionID", connectionId,
			"PortID", portID,
			"Error", err,
		)
		return nil
	}

	// get zone info
	zoneInfo, found := im.keeper.GetRegisteredZoneInfo(ctx, chainId)
	if !found {
		ctx.Logger().Error(fmt.Sprintf("expected to find zone info for %v", chainId))
		return nil
	}

	// get interchain account address
	address, found := im.keeper.ICAControllerKeeper.GetInterchainAccountAddress(ctx, connectionId, portID)
	if !found {
		ctx.Logger().Error(fmt.Sprintf("expected to find an address for %s/%s", connectionId, portID))
		return nil
	}

	ctx.Logger().Info("Found matching address", "chain", zoneInfo.ChainId, "address", address, "port", portID)
	portParts := strings.Split(portID, ".")

	switch {
	// deposit address
	case len(portParts) == 2 && portParts[1] == "deposit":

		// refactor: register DepositAddress

		zoneInfo.DepositAddress = &types.ICAAccount{Address: address, Balance: sdk.Coins{}, DelegatedBalance: sdk.NewCoin(zoneInfo.BaseDenom, sdk.ZeroInt()), PortName: portID}

		balanceQuery := bankTypes.QueryAllBalancesRequest{Address: address}
		bz, err := im.keeper.GetCodec().Marshal(&balanceQuery)
		if err != nil {
			return err
		}

		im.keeper.ICQKeeper.MakeRequest(
			ctx,
			connectionId,
			chainId,
			"cosmos.bank.v1beta1.Query/AllBalances",
			bz,
			sdk.NewInt(int64(im.keeper.GetParam(ctx, types.KeyDepositInterval))),
			types.ModuleName,
			"allbalances",
			0,
		)

	// withdrawal address
	case len(portParts) == 2 && portParts[1] == "withdrawal":

		// refactor: register WithdrawalAddress

		zoneInfo.WithdrawalAddress = &types.ICAAccount{Address: address, Balance: sdk.Coins{}, DelegatedBalance: sdk.NewCoin(zoneInfo.BaseDenom, sdk.ZeroInt()), PortName: portID}

		for _, da := range zoneInfo.GetDelegationAccounts() {
			msg := distrTypes.MsgSetWithdrawAddress{DelegatorAddress: da.Address, WithdrawAddress: address}
			im.keeper.SubmitTx(ctx, []sdk.Msg{&msg}, da, "")
		}

	// delegation addresses
	case len(portParts) == 3 && portParts[1] == "delegate":

		// refactor: register DelegationAddresses

		delegationAccounts := zoneInfo.GetDelegationAccounts()
		// check for duplicate address
		for _, existing := range delegationAccounts {
			if existing.Address == address {
				ctx.Logger().Error("unexpectedly found existing address: " + address)
				return nil
			}
		}
		account := &types.ICAAccount{Address: address, Balance: sdk.Coins{}, DelegatedBalance: sdk.NewCoin(zoneInfo.BaseDenom, sdk.ZeroInt()), PortName: portID}
		// append delegation account address
		zoneInfo.DelegationAddresses = append(delegationAccounts, account)

		// set withdrawal address if, and only if withdrawal address is already set
		if zoneInfo.WithdrawalAddress != nil {
			msg := distrTypes.MsgSetWithdrawAddress{DelegatorAddress: address, WithdrawAddress: zoneInfo.WithdrawalAddress.String()}
			im.keeper.SubmitTx(ctx, []sdk.Msg{&msg}, account, "")
		}

		// var cb keeper.Callback = func(k keeper.Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
		// 	zone, found := k.GetRegisteredZoneInfo(ctx, query.GetChainId())
		// 	if !found {
		// 		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
		// 	}
		// 	// unmarshal request
		// 	balanceQuery := bankTypes.QueryAllBalancesRequest{}
		// 	err := k.GetCodec().Unmarshal(query.Request, &balanceQuery)
		// 	if err == nil {
		// 		return err
		// 	}

		// 	return k.SetAccountBalance(ctx, zone, balanceQuery.Address, args)
		// }

		// balanceQuery := bankTypes.QueryAllBalancesRequest{Address: address}
		// bz, err := im.keeper.GetCodec().Marshal(&balanceQuery)
		// if err != nil {
		// 	return err
		// }

		// im.keeper.ICQKeeper.MakeRequest(
		// 	ctx,
		// 	connectionId,
		// 	chainId,
		// 	"cosmos.bank.v1beta1.Query/AllBalances",
		// 	bz,
		// 	sdk.NewInt(int64(im.keeper.GetParam(ctx, types.KeyDelegateInterval))),
		// 	types.ModuleName,
		// 	cb,
		// )

	// performance address
	case len(portParts) == 2 && portParts[1] == "performance":

		if err := im.registerPerformanceAddress(ctx, portID, connectionId, chainId, address, &zoneInfo); err != nil {
			return err
		}

	default:
		ctx.Logger().Error("unexpected channel on portID: " + portID)

	}
	im.keeper.SetRegisteredZone(ctx, zoneInfo)
	return nil
}

func (im IBCModule) registerPerformanceAddress(
	ctx sdk.Context,
	portId,
	connectionId,
	chainId,
	address string,
	zone *types.RegisteredZone,
) error {
	zone.PerformanceAddress = &types.ICAAccount{Address: address, Balance: sdk.Coins{}, DelegatedBalance: sdk.NewCoin(zone.BaseDenom, sdk.ZeroInt()), PortName: portId}

	// set withdrawal address if, and only if withdrawal address is already set
	if zone.WithdrawalAddress != nil {
		msg := distrTypes.MsgSetWithdrawAddress{DelegatorAddress: address, WithdrawAddress: zone.WithdrawalAddress.String()}
		im.keeper.SubmitTx(ctx, []sdk.Msg{&msg}, zone.PerformanceAddress, "")
	}

	return im.keeper.EmitPerformanceBalanceQuery(ctx, zone)
}

// OnChanOpenConfirm implements the IBCModule interface
func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

// OnChanCloseInit implements the IBCModule interface
func (im IBCModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

// OnChanCloseConfirm implements the IBCModule interface
func (im IBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

// OnRecvPacket implements the IBCModule interface. A successful acknowledgement
// is returned if the packet data is successfully decoded and the receive application
// logic returns without error.
func (im IBCModule) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	return channeltypes.NewErrorAcknowledgement(fmt.Errorf("cannot receive packet via interchain accounts authentication module"))
}

// OnAcknowledgementPacket implements the IBCModule interface
func (im IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	connectionId, _, err := im.keeper.IBCKeeper.ChannelKeeper.GetChannelConnection(ctx, packet.SourcePort, packet.SourceChannel)
	if err != nil {
		err = fmt.Errorf("packet connection not found: %w", err)
		ctx.Logger().Error(err.Error())
		return err
	}
	ctx = ctx.WithContext(context.WithValue(ctx.Context(), "connectionId", connectionId))

	return im.keeper.HandleAcknowledgement(ctx, packet, acknowledgement)
}

// OnTimeoutPacket implements the IBCModule interface.
func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	return im.keeper.HandleTimeout(ctx, packet)
}

// NegotiateAppVersion implements the IBCModule interface
func (im IBCModule) NegotiateAppVersion(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionID string,
	portID string,
	counterparty channeltypes.Counterparty,
	proposedVersion string,
) (string, error) {
	return "", nil
}
