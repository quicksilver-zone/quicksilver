package app

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cometbft/cometbft/abci/types"
)

var (
	nodeID      string
	nodeIDMutex sync.Mutex
)

func (app *Quicksilver) InitNodeID() (string, error) {
	nodeIDMutex.Lock()
	defer nodeIDMutex.Unlock()
	if nodeID == "" {
		query := tmservice.GetNodeInfoRequest{}
		queryRes := app.Query(types.RequestQuery{Data: app.appCodec.MustMarshal(&query), Path: "/cosmos.base.tendermint.v1beta1.Service/GetNodeInfo", Height: 0, Prove: false})
		nodeInfo := queryRes.GetValue()
		result := tmservice.GetNodeInfoResponse{}
		err := app.appCodec.Unmarshal(nodeInfo, &result)
		if err != nil {
			return "", err
		}
		if result.DefaultNodeInfo == nil {
			return "", errors.New("node info is nil")
		}
		nodeID = result.DefaultNodeInfo.DefaultNodeID
	}
	return nodeID, nil
}

func (app *Quicksilver) ShipMetrics(ctx sdk.Context) {
	metricsFamilies, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		ctx.Logger().Error("Error gathering metrics", "error", err)
	}

	buf := &bytes.Buffer{}
	defer buf.Reset()

	e := expfmt.NewEncoder(buf, expfmt.NewFormat(expfmt.TypeTextPlain))
	for _, mf := range metricsFamilies {
		if err := e.Encode(mf); err != nil {
			ctx.Logger().Error("Error encoding metrics", "error", err)
			return
		}
	}

	nodeID, err := app.InitNodeID()
	if err != nil {
		ctx.Logger().Error("Error getting node ID", "error", err)
		return
	}
	c, cancel := context.WithTimeout(ctx.Context(), time.Second*10)
	defer cancel()
	req, err := http.NewRequestWithContext(c, "POST", fmt.Sprintf("%s/job/%s/instance/%s", app.metricsURL, ctx.ChainID(), nodeID), bytes.NewReader(buf.Bytes()))
	if err != nil {
		ctx.Logger().Error("Error creating metrics request", "error", err)
		return
	}
	req.Header.Set("Content-Type", string(expfmt.NewFormat(expfmt.TypeTextPlain)))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ctx.Logger().Error("Error sending metrics", "error", err)
		return
	}
	defer resp.Body.Close()
}
