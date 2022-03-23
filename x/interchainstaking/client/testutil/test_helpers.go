package testutil

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/client/cli"
)

var (
	commonArgs = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
	}
)

func MsgRegisterZoneExec(clientCtx client.Context, moniker, nodeId, chainId, denom, from string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		moniker,
		nodeId,
		chainId,
		denom,
		denom,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}

	args = append(args, extraArgs...)
	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, cli.GetRegisterZoneTxCmd(), args)
}
