package types

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"

	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/logger"
)

const (
	DataTypeProtocolData = 0
	DataTypeZone         = 1
)

type Data[T any] struct {
	Data []T
}

type Zone[T any] struct {
	Zones      []T
	Stats      any
	Pagination any
}

func NewCacheManager() CacheManager {
	return CacheManager{Data: make(map[string]CacheManagerElementI, 0)}
}

func GetCache[T prewards.ConnectionProtocolData | prewards.OsmosisParamsProtocolData | prewards.OsmosisPoolProtocolData | prewards.OsmosisClPoolProtocolData | prewards.LiquidAllowedDenomProtocolData | prewards.UmeeParamsProtocolData | icstypes.Zone](ctx context.Context, mgr *CacheManager) ([]T, error) {
	cache, ok := mgr.Data[new(Cache[T]).Type()].(*Cache[T])
	if !ok {
		return nil, fmt.Errorf("cache not found for type %T", new(T))
	}
	value := cache.Get(ctx)
	return value, nil
}

func AddMocks[T prewards.ConnectionProtocolData | prewards.OsmosisParamsProtocolData | prewards.OsmosisPoolProtocolData | prewards.OsmosisClPoolProtocolData | prewards.LiquidAllowedDenomProtocolData | prewards.UmeeParamsProtocolData | icstypes.Zone](ctx context.Context, mgr *CacheManager, mocks []T) {
	cache, _ := mgr.Data[new(Cache[T]).Type()].(*Cache[T])
	cache.SetMock(mocks)
}

type CacheManager struct {
	Data map[string]CacheManagerElementI
}

type CacheManagerElementI interface {
	Init(ctx context.Context, url string, dataType int, updateTime time.Duration) error
	Fetch(ctx context.Context)
	Type() string
}

func (m *CacheManager) Add(ctx context.Context, element CacheManagerElementI, url string, dataType int, updateTime time.Duration) error {
	m.Data[element.Type()] = element
	return m.Data[element.Type()].Init(ctx, url, dataType, updateTime)
}

// GetConnections implements CacheManagerInterface
func (m *CacheManager) GetConnections(ctx context.Context) ([]prewards.ConnectionProtocolData, error) {
	return GetCache[prewards.ConnectionProtocolData](ctx, m)
}

// GetOsmosisParams implements CacheManagerInterface
func (m *CacheManager) GetOsmosisParams(ctx context.Context) ([]prewards.OsmosisParamsProtocolData, error) {
	return GetCache[prewards.OsmosisParamsProtocolData](ctx, m)
}

// GetOsmosisPools implements CacheManagerInterface
func (m *CacheManager) GetOsmosisPools(ctx context.Context) ([]prewards.OsmosisPoolProtocolData, error) {
	return GetCache[prewards.OsmosisPoolProtocolData](ctx, m)
}

// GetOsmosisClPools implements CacheManagerInterface
func (m *CacheManager) GetOsmosisClPools(ctx context.Context) ([]prewards.OsmosisClPoolProtocolData, error) {
	return GetCache[prewards.OsmosisClPoolProtocolData](ctx, m)
}

// GetLiquidAllowedDenoms implements CacheManagerInterface
func (m *CacheManager) GetLiquidAllowedDenoms(ctx context.Context) ([]prewards.LiquidAllowedDenomProtocolData, error) {
	return GetCache[prewards.LiquidAllowedDenomProtocolData](ctx, m)
}

// GetUmeeParams implements CacheManagerInterface
func (m *CacheManager) GetUmeeParams(ctx context.Context) ([]prewards.UmeeParamsProtocolData, error) {
	return GetCache[prewards.UmeeParamsProtocolData](ctx, m)
}

// GetZones implements CacheManagerInterface
func (m *CacheManager) GetZones(ctx context.Context) ([]icstypes.Zone, error) {
	return GetCache[icstypes.Zone](ctx, m)
}

// AddMocks implements CacheManagerInterface
func (m *CacheManager) AddMocks(ctx context.Context, mocks interface{}) error {
	// This is a simplified implementation - in practice you'd need to handle different types
	return nil
}

type CacheI[T any] interface {
	Init(ctx context.Context, url string, dataType int, updateTime time.Duration) error
	Fetch(ctx context.Context)
	Get(ctx context.Context) []T
}

var (
	_ CacheI[any]          = &Cache[any]{}
	_ CacheManagerElementI = &Cache[any]{}
)

type Cache[T any] struct {
	dataType    int
	url         string
	cache       []T
	lastUpdated time.Time
	duration    time.Duration
	mockData    []T
}

func (c *Cache[T]) Type() string {
	a := new(T)
	return strings.ReplaceAll(reflect.TypeOf(*a).String(), "types.", "")
}

func (c *Cache[T]) unmarshal(responseData []byte) []T {
	switch c.dataType {
	case DataTypeProtocolData:
		data := Data[T]{}

		err := json.Unmarshal(responseData, &data)
		if err != nil {
			panic(err)
		}
		return data.Data
	case DataTypeZone:
		data := Zone[T]{}

		err := json.Unmarshal(responseData, &data)
		if err != nil {
			panic(err)
		}
		return data.Zones
	}
	return nil
}

func (c *Cache[T]) Init(ctx context.Context, url string, dataType int, updateInterval time.Duration) error {
	c.url = url
	c.duration = updateInterval
	c.dataType = dataType
	c.Fetch(ctx)
	return nil
}

func (c *Cache[T]) SetMock(mocks []T) {
	c.mockData = mocks
}

func (c *Cache[T]) read(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, http.NoBody)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return io.ReadAll(response.Body)
}

func (c *Cache[T]) Fetch(ctx context.Context) {
	log := logger.FromContext(ctx)
	log.Debug("Fetching and caching data", "url", c.url)

	responseData, err := c.read(ctx)
	if err != nil {
		log.Error("Failed to read cache data", "error", err, "url", c.url)
		return
	}

	c.cache = c.unmarshal(responseData)
	c.lastUpdated = time.Now()
}

func (c Cache[T]) Get(ctx context.Context) []T {
	if time.Now().After(c.lastUpdated.Add(c.duration)) {
		c.Fetch(ctx)
	}
	return append(c.cache, c.mockData...)
}
