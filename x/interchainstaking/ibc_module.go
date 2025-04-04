package interchainstaking

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"

	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"

	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/keeper"
)

var _ porttypes.IBCModule = IBCModule{}

// IBCModule implements the ICS26 interface for interchain accounts controller chains.
type IBCModule struct {
	keeper    *keeper.Keeper
	icaKeeper *icacontrollerkeeper.Keeper
}

// NewIBCModule creates a new IBCModule given the keeper.
func NewIBCModule(k *keeper.Keeper) IBCModule {
	return IBCModule{
		keeper: k,
	}
}

// OnChanOpenInit implements the IBCModule interface.
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
	return "", nil
}

// OnChanOpenTry implements the IBCModule interface.
func (IBCModule) OnChanOpenTry(
	ctx sdk.Context,
	_ channeltypes.Order,
	_ []string,
	_ string,
	_ string,
	_ *capabilitytypes.Capability,
	_ channeltypes.Counterparty,
	version string,
) (string, error) {
	return version, nil
}

// OnChanOpenAck implements the IBCModule interface.
func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	// get connection from port
	connectionID, _, err := im.keeper.IBCKeeper.ChannelKeeper.GetChannelConnection(ctx, portID, channelID)
	if err != nil {
		return err
	}

	return im.keeper.HandleChannelOpenAck(ctx, portID, connectionID)
}

// OnChanOpenConfirm implements the IBCModule interface.
func (im IBCModule) OnChanOpenConfirm(
	_ sdk.Context,
	_,
	_ string,
) error {
	return nil
}

// OnChanCloseInit implements the IBCModule interface.
func (im IBCModule) OnChanCloseInit(
	_ sdk.Context,
	_,
	_ string,
) error {
	return nil
}

// OnChanCloseConfirm implements the IBCModule interface.
func (im IBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return im.icaKeeper.OnChanCloseConfirm(ctx, portID, channelID)
}

// OnRecvPacket implements the IBCModule interface. A successful acknowledgement
// is returned if the packet data is successfully decoded and the receive application
// logic returns without error.
func (im IBCModule) OnRecvPacket(
	_ sdk.Context,
	_ channeltypes.Packet,
	_ sdk.AccAddress,
) ibcexported.Acknowledgement {
	return channeltypes.NewErrorAcknowledgement(errors.New("cannot receive packet via interchainstaking authentication module"))
}

// OnAcknowledgementPacket implements the IBCModule interface.
func (im IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	_ sdk.AccAddress,
) error {
	connectionID, _, err := im.keeper.IBCKeeper.ChannelKeeper.GetChannelConnection(ctx, packet.SourcePort, packet.SourceChannel)
	if err != nil {
		err = fmt.Errorf("packet connection not found: %w", err)
		ctx.Logger().Error(err.Error())
		return err
	}
	err = im.keeper.HandleAcknowledgement(ctx, packet, acknowledgement, connectionID)
	if err != nil {
		im.keeper.Logger(ctx).Error("CALLBACK ERROR:", "error", err.Error())
	}
	return err
}

// OnTimeoutPacket implements the IBCModule interface.
func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	_ sdk.AccAddress,
) error {
	return im.keeper.HandleTimeout(ctx, packet)
}

// NegotiateAppVersion implements the IBCModule interface.
func (IBCModule) NegotiateAppVersion(
	_ sdk.Context,
	_ channeltypes.Order,
	_ string,
	_ string,
	_ channeltypes.Counterparty,
	proposedVersion string,
) (string, error) {
	return proposedVersion, nil
}
