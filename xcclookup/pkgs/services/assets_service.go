package services

import (
	"context"
	"errors"
	"fmt"
	"sync"

	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"

	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types"
)

// AssetsService handles assets-related operations
type AssetsService struct {
	cfg           types.Config
	cacheManager  types.CacheManagerInterface
	claimsService types.ClaimsServiceInterface
	heights       map[string]int64
}

// NewAssetsService creates a new assets service
func NewAssetsService(
	cfg types.Config,
	cacheManager types.CacheManagerInterface,
	claimsService types.ClaimsServiceInterface,
	heights map[string]int64,
) *AssetsService {
	return &AssetsService{
		cfg:           cfg,
		cacheManager:  cacheManager,
		claimsService: claimsService,
		heights:       heights,
	}
}

// GetAssets retrieves assets for a given address
func (s *AssetsService) GetAssets(ctx context.Context, address string) (*types.Response, map[string]error) {
	errs := make(map[string]error)
	errsMutex := sync.Mutex{}
	response := &types.Response{
		Messages: make([]prewards.MsgSubmitClaim, 0),
		Assets:   make(map[string][]types.Asset),
	}

	wg := sync.WaitGroup{}

	// Get connections
	unfilteredConnections, err := s.cacheManager.GetConnections(ctx)
	if err != nil {
		errs["Connections"] = err
		return response, errs
	}

	var connections []prewards.ConnectionProtocolData
	for _, ufc := range unfilteredConnections {
		if ufc.LastEpoch > 0 {
			connections = append(connections, ufc)
		}
	}

	mappedAddresses, err := types.GetMappedAddresses(ctx, address, unfilteredConnections, &s.cfg)
	if err != nil {
		errs["MappedAddresses"] = err
	}

	// Process Osmosis claims
	s.processOsmosisClaims(ctx, address, mappedAddresses, response, errs, &errsMutex, &wg)

	// Process Umee claims
	s.processUmeeClaims(ctx, address, mappedAddresses, response, errs, &errsMutex, &wg)

	// Process Membrane claims
	s.processMembraneClaims(ctx, address, mappedAddresses, response, errs, &errsMutex, &wg)

	// Process Liquid claims
	s.processLiquidClaims(ctx, address, mappedAddresses, connections, response, errs, &errsMutex, &wg)

	wg.Wait()
	return response, errs
}

func (s *AssetsService) processOsmosisClaims(
	ctx context.Context,
	address string,
	mappedAddresses map[string]string,
	response *types.Response,
	errs map[string]error,
	errsMutex *sync.Mutex,
	wg *sync.WaitGroup,
) {
	osmosisParamsCache, err := s.cacheManager.GetOsmosisParams(ctx)
	if err != nil {
		errsMutex.Lock()
		errs["OsmosisParams"] = err
		errsMutex.Unlock()
		return
	}
	if len(osmosisParamsCache) == 0 {
		errsMutex.Lock()
		errs["OsmosisConfig"] = errors.New("osmosis params not set")
		errsMutex.Unlock()
		return
	}

	chain := osmosisParamsCache[0].ChainID

	wg.Add(1)
	go func() {
		defer wg.Done()
		result, err := s.claimsService.OsmosisClaim(ctx, address, address, chain, s.heights[chain])
		if err != nil {
			errsMutex.Lock()
			errs["OsmosisClaim"] = err
			errsMutex.Unlock()
		}
		if result.Err != nil {
			errsMutex.Lock()
			errs["OsmosisClaim"] = result.Err
			errsMutex.Unlock()
		}
		if result.OsmosisPool.Err != nil {
			errsMutex.Lock()
			errs["OsmosisPoolClaim"] = result.OsmosisPool.Err
			errsMutex.Unlock()
		}
		if result.OsmosisClPool.Err != nil {
			errsMutex.Lock()
			errs["OsmosisClPoolClaim"] = result.OsmosisClPool.Err
			errsMutex.Unlock()
		}
		if result.OsmosisPool.Msg != nil {
			response.Update(ctx, result.OsmosisPool.Msg, result.OsmosisPool.Assets, "osmosispool")
		}
		if result.OsmosisClPool.Msg != nil {
			response.Update(ctx, result.OsmosisClPool.Msg, result.OsmosisClPool.Assets, "osmosisclpool")
		}
	}()

	if mappedAddress, ok := mappedAddresses[chain]; ok {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := s.claimsService.OsmosisClaim(ctx, mappedAddress, address, chain, s.heights[chain])
			if err != nil {
				errsMutex.Lock()
				errs["OsmosisClaim"] = err
				errsMutex.Unlock()
			}
			if result.Err != nil {
				errsMutex.Lock()
				errs["OsmosisClaim"] = result.Err
				errsMutex.Unlock()
			}
			if result.OsmosisPool.Err != nil {
				errsMutex.Lock()
				errs["OsmosisPoolClaim"] = result.OsmosisPool.Err
				errsMutex.Unlock()
			}
			if result.OsmosisClPool.Err != nil {
				errsMutex.Lock()
				errs["OsmosisClPoolClaim"] = result.OsmosisClPool.Err
				errsMutex.Unlock()
			}
			if result.OsmosisPool.Msg != nil {
				response.Update(ctx, result.OsmosisPool.Msg, result.OsmosisPool.Assets, "osmosispool")
			}
			if result.OsmosisClPool.Msg != nil {
				response.Update(ctx, result.OsmosisClPool.Msg, result.OsmosisClPool.Assets, "osmosisclpool")
			}
		}()
	}
}

func (s *AssetsService) processUmeeClaims(
	ctx context.Context,
	address string,
	mappedAddresses map[string]string,
	response *types.Response,
	errs map[string]error,
	errsMutex *sync.Mutex,
	wg *sync.WaitGroup,
) {
	umeeParamsCache, err := s.cacheManager.GetUmeeParams(ctx)
	if err != nil {
		errsMutex.Lock()
		errs["UmeeParams"] = err
		errsMutex.Unlock()
		return
	}
	if len(umeeParamsCache) == 0 {
		errsMutex.Lock()
		errs["UmeeConfig"] = errors.New("umee params not set")
		errsMutex.Unlock()
		return
	}

	chain := umeeParamsCache[0].ChainID

	wg.Add(1)
	go func() {
		defer wg.Done()
		messages, assets, err := s.claimsService.UmeeClaim(ctx, address, address, chain, s.heights[chain])
		if err != nil {
			errsMutex.Lock()
			errs["UmeeClaim"] = err
			errsMutex.Unlock()
			return
		}
		if messages != nil {
			response.Update(ctx, messages, assets, "umee")
		}
	}()

	if mappedAddress, ok := mappedAddresses[chain]; ok {
		wg.Add(1)
		go func() {
			defer wg.Done()
			messages, assets, err := s.claimsService.UmeeClaim(ctx, mappedAddress, address, chain, s.heights[chain])
			if err != nil {
				errsMutex.Lock()
				errs["UmeeClaim"] = err
				errsMutex.Unlock()
			}
			if messages != nil {
				response.Update(ctx, messages, assets, "umee")
			}
		}()
	}
}

func (s *AssetsService) processMembraneClaims(
	ctx context.Context,
	address string,
	mappedAddresses map[string]string,
	response *types.Response,
	errs map[string]error,
	errsMutex *sync.Mutex,
	wg *sync.WaitGroup,
) {
	membraneParamsCache, err := s.cacheManager.GetMembraneParams(ctx)
	if err != nil {
		errsMutex.Lock()
		errs["MembraneParams"] = err
		errsMutex.Unlock()
		return
	}
	if len(membraneParamsCache) == 0 {
		errsMutex.Lock()
		errs["MembraneConfig"] = errors.New("membrane params not set")
		errsMutex.Unlock()
		return
	}

	// Get osmosis params to determine the chain
	osmosisParamsCache, err := s.cacheManager.GetOsmosisParams(ctx)
	if err != nil {
		errsMutex.Lock()
		errs["OsmosisParams"] = err
		errsMutex.Unlock()
		return
	}
	if len(osmosisParamsCache) == 0 {
		errsMutex.Lock()
		errs["OsmosisConfig"] = errors.New("osmosis params not set")
		errsMutex.Unlock()
		return
	}

	chain := osmosisParamsCache[0].ChainID

	wg.Add(1)
	go func() {
		defer wg.Done()
		messages, assets, err := s.claimsService.MembraneClaim(ctx, address, address, chain, s.heights[chain])
		if err != nil {
			errsMutex.Lock()
			errs["MembraneClaim"] = err
			errsMutex.Unlock()
			return
		}
		if messages != nil {
			response.Update(ctx, messages, assets, "membrane")
		}
	}()

	if mappedAddress, ok := mappedAddresses[chain]; ok {
		wg.Add(1)
		go func() {
			defer wg.Done()
			messages, assets, err := s.claimsService.MembraneClaim(ctx, mappedAddress, address, chain, s.heights[chain])
			if err != nil {
				errsMutex.Lock()
				errs["MembraneClaim"] = err
				errsMutex.Unlock()
			}
			if messages != nil {
				response.Update(ctx, messages, assets, "membrane")
			}
		}()
	}
}

func (s *AssetsService) processLiquidClaims(
	ctx context.Context,
	address string,
	mappedAddresses map[string]string,
	connections []prewards.ConnectionProtocolData,
	response *types.Response,
	errs map[string]error,
	errsMutex *sync.Mutex,
	wg *sync.WaitGroup,
) {
	for _, con := range connections {
		wg.Add(1)
		go func(con prewards.ConnectionProtocolData) {
			defer wg.Done()
			messages, assets, err := s.claimsService.LiquidClaim(ctx, address, address, con, s.heights[con.ChainID])
			if err != nil {
				errsMutex.Lock()
				errs[fmt.Sprintf("LiquidClaim:%s", con.ChainID)] = err
				errsMutex.Unlock()
				return
			}
			response.Update(ctx, messages, assets, "liquid")
		}(con)

		if mappedAddress, ok := mappedAddresses[con.ChainID]; ok {
			wg.Add(1)
			go func(con prewards.ConnectionProtocolData) {
				defer wg.Done()
				messages, assets, err := s.claimsService.LiquidClaim(ctx, mappedAddress, address, con, s.heights[con.ChainID])
				if err != nil {
					errsMutex.Lock()
					errs[fmt.Sprintf("LiquidClaim:%s", con.ChainID)] = err
					errsMutex.Unlock()
					return
				}
				response.Update(ctx, messages, assets, "liquid")
			}(con)
		}
	}
}
