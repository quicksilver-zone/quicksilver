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

func GetCache[T prewards.ConnectionProtocolData | prewards.OsmosisParamsProtocolData | prewards.OsmosisPoolProtocolData | prewards.OsmosisClPoolProtocolData | prewards.LiquidAllowedDenomProtocolData | prewards.UmeeParamsProtocolData | icstypes.Zone](ctx context.Context, mgr *CacheManager) ([]T, error) {
	cache, _ := mgr.Data[new(Cache[T]).Type()].(*Cache[T])
	value, err := cache.Get(ctx)
	if err != nil {
		return nil, err
	}
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
	Fetch(ctx context.Context) error
	Type() string
}

func (m *CacheManager) Add(ctx context.Context, element CacheManagerElementI, url string, dataType int, updateTime time.Duration) error {
	m.Data[element.Type()] = element
	err := m.Data[element.Type()].Init(ctx, url, dataType, updateTime)
	if err != nil {
		return err
	}
	return nil
}

type CacheI[T any] interface {
	Init(ctx context.Context, url string, dataType int, updateTime time.Duration) error
	Fetch(ctx context.Context) error
	Get(ctx context.Context) ([]T, error)
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

func (c *Cache[T]) unmarshal(responseData []byte) ([]T, error) {
	switch c.dataType {
	case DataTypeProtocolData:
		data := Data[T]{}

		err := json.Unmarshal(responseData, &data)
		if err != nil {
			return nil, err
		}
		return data.Data, nil
	case DataTypeZone:
		data := Zone[T]{}

		err := json.Unmarshal(responseData, &data)
		if err != nil {
			return nil, err
		}
		return data.Zones, nil
	}
	return nil, fmt.Errorf("invalid data type: %d", c.dataType)
}

func (c *Cache[T]) Init(ctx context.Context, url string, dataType int, updateInterval time.Duration) error {
	c.url = url
	c.duration = updateInterval
	c.dataType = dataType
	return c.Fetch(ctx)
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

func (c *Cache[T]) Fetch(ctx context.Context) error {
	fmt.Println("Fetching and caching " + c.url)

	responseData, err := c.read(ctx)
	if err != nil {
		return err
	}

	c.cache, err = c.unmarshal(responseData)
	if err != nil {
		return err
	}
	c.lastUpdated = time.Now()
	return nil
}

func (c *Cache[T]) Get(ctx context.Context) ([]T, error) {
	if time.Now().After(c.lastUpdated.Add(c.duration)) {
		err := c.Fetch(ctx)
		if err != nil {
			return nil, err
		}
	}
	return append(c.cache, c.mockData...), nil
}
