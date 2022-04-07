package interchainstaking

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"

	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v3/modules/core/exported"
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
) error {
	return im.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID))
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
	// TODO: is there re-entrancy risk here?
	connectionId, err := im.keeper.GetConnectionForPort(ctx, portID)
	if err != nil {
		ctx.Logger().Error("Unable to get connection for port " + portID)
	}
	address, found := im.keeper.ICAControllerKeeper.GetInterchainAccountAddress(ctx, connectionId, portID)
	if !found {
		ctx.Logger().Error(fmt.Sprintf("Expected to find an address for %s/%s", connectionId, portID))
		return nil
	}
	im.keeper.IterateRegisteredZones(ctx, func(index int64, zoneInfo types.RegisteredZone) (stop bool) {
		if zoneInfo.GetConnectionId() == connectionId {
			if !strings.Contains(portID, zoneInfo.ChainId) {
				im.keeper.DeleteRegisteredZone(ctx, zoneInfo.ChainId)
				return false
			}
			ctx.Logger().Info("Found matching address", "chain", zoneInfo.ChainId, "address", address, "port", portID)
			portParts := strings.Split(portID, ".")
			if len(portParts) == 2 && portParts[1] == "deposit" {
				zoneInfo.DepositAddress = &types.ICAAccount{Address: address, Balance: sdk.Coins{}, DelegatedBalance: sdk.Coin{}, PortName: portID}
				balanceQuery := im.keeper.ICQKeeper.NewPeriodicQuery(ctx, connectionId, zoneInfo.ChainId, "cosmos.bank.v1beta1.Query/AllBalances", map[string]string{"address": address}, sdk.NewInt(int64(im.keeper.GetParam(ctx, types.KeyDepositInterval))))
				im.keeper.ICQKeeper.SetPeriodicQuery(ctx, *balanceQuery)
			} else if len(portParts) == 3 && portParts[1] == "delegate" {
				for _, existing := range zoneInfo.DelegationAddresses {
					if existing.Address == address {
						ctx.Logger().Error("unexpectedly found existing address: " + address)
						return false
					}
				}
				zoneInfo.DelegationAddresses = append(zoneInfo.DelegationAddresses, &types.ICAAccount{Address: address, Balance: sdk.Coins{}, DelegatedBalance: sdk.Coin{}, PortName: portID})
				balanceQuery := im.keeper.ICQKeeper.NewPeriodicQuery(ctx, connectionId, zoneInfo.ChainId, "cosmos.bank.v1beta1.Query/AllBalances", map[string]string{"address": address}, sdk.NewInt(int64(im.keeper.GetParam(ctx, types.KeyDelegateInterval))))
				im.keeper.ICQKeeper.SetPeriodicQuery(ctx, *balanceQuery)
				rewardsQuery := im.keeper.ICQKeeper.NewPeriodicQuery(ctx, connectionId, zoneInfo.ChainId, "cosmos.distribution.v1beta1.Query/DelegationTotalRewards", map[string]string{"delegator": address}, sdk.NewInt(int64(im.keeper.GetParam(ctx, types.KeyDelegateInterval))))
				im.keeper.ICQKeeper.SetPeriodicQuery(ctx, *rewardsQuery)
				delegationQuery := im.keeper.ICQKeeper.NewPeriodicQuery(ctx, connectionId, zoneInfo.ChainId, "cosmos.staking.v1beta1.Query/DelegatorDelegations", map[string]string{"address": address}, sdk.NewInt(int64(im.keeper.GetParam(ctx, types.KeyDelegationsInterval)))) // this can probably be less frequent, because we manage delegations ourselves.
				im.keeper.ICQKeeper.SetPeriodicQuery(ctx, *delegationQuery)
			} else {
				ctx.Logger().Error("unexpected channel on portID: " + portID)
				return false
			}
			// save zone
			im.keeper.SetRegisteredZone(ctx, zoneInfo)
			return true
		}
		return false
	})
	return nil
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
	return channeltypes.NewErrorAcknowledgement("cannot receive packet via interchain accounts authentication module")
}

// OnAcknowledgementPacket implements the IBCModule interface
func (im IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
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
