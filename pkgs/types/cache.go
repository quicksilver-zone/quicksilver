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

func GetCache[T prewards.ConnectionProtocolData | prewards.OsmosisParamsProtocolData | prewards.OsmosisPoolProtocolData | prewards.OsmosisClPoolProtocolData | prewards.LiquidAllowedDenomProtocolData | prewards.UmeeParamsProtocolData | icstypes.Zone](ctx context.Context, mgr *CacheManager) []T {
	cache, _ := mgr.Data[new(Cache[T]).Type()].(*Cache[T])
	value := cache.Get(ctx)
	return value
}

func AddMocks[T prewards.ConnectionProtocolData | prewards.OsmosisParamsProtocolData | prewards.OsmosisPoolProtocolData | prewards.OsmosisClPoolProtocolData | prewards.LiquidAllowedDenomProtocolData | prewards.UmeeParamsProtocolData | icstypes.Zone](ctx context.Context, mgr *CacheManager, mocks []T) {
	cache, _ := mgr.Data[new(Cache[T]).Type()].(*Cache[T])
	cache.SetMock(mocks)
}

type CacheManager struct {
	Data map[string]CacheManagerElementI
}

type CacheManagerElementI interface {
	Init(ctx context.Context, url string, dataType int, updateTime time.Duration)
	Fetch(ctx context.Context)
	Type() string
}

func (m *CacheManager) Add(ctx context.Context, element CacheManagerElementI, url string, dataType int, updateTime time.Duration) {
	m.Data[element.Type()] = element
	m.Data[element.Type()].Init(ctx, url, dataType, updateTime)
}

type CacheI[T any] interface {
	Init(ctx context.Context, url string, dataType int, updateTime time.Duration)
	Fetch(ctx context.Context)
	Get(ctx context.Context) []T
}

var _ CacheI[any] = &Cache[any]{}
var _ CacheManagerElementI = &Cache[any]{}

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
	return strings.Replace(reflect.TypeOf(*a).String(), "types.", "", -1)
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

func (c *Cache[T]) Init(ctx context.Context, url string, dataType int, updateInterval time.Duration) {
	c.url = url
	c.duration = updateInterval
	c.dataType = dataType
	c.Fetch(ctx)
}

func (c *Cache[T]) SetMock(mocks []T) {
	c.mockData = mocks
}

func (c *Cache[T]) read(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, http.NoBody)
	if err != nil {
		panic(err)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	return io.ReadAll(response.Body)
}

func (c *Cache[T]) Fetch(ctx context.Context) {
	fmt.Println("Fetching and caching " + c.url)

	responseData, err := c.read(ctx)
	if err != nil {
		panic(err)
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
