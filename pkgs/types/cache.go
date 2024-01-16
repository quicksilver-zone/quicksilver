package types

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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

type CacheManagerI[T any] interface {
	Init(ctx context.Context, url string, dataType int, updateTime time.Duration)
	Fetch(ctx context.Context)
	Get(ctx context.Context) []T
}

var _ CacheManagerI[any] = &CacheManager[any]{}

type CacheManager[T any] struct {
	dataType    int
	url         string
	cache       []T
	lastUpdated time.Time
	duration    time.Duration
}

func (c *CacheManager[T]) unmarshal(responseData []byte) []T {
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

func (c *CacheManager[T]) Init(ctx context.Context, url string, dataType int, updateInterval time.Duration) {
	c.url = url
	c.duration = updateInterval
	c.dataType = dataType
	c.Fetch(ctx)
}

func (c *CacheManager[T]) read(ctx context.Context) ([]byte, error) {
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

func (c *CacheManager[T]) Fetch(ctx context.Context) {
	fmt.Println("Fetching and caching " + c.url)

	responseData, err := c.read(ctx)
	if err != nil {
		panic(err)
	}

	c.cache = c.unmarshal(responseData)
	c.lastUpdated = time.Now()
}

func (c CacheManager[T]) Get(ctx context.Context) []T {
	if time.Now().After(c.lastUpdated.Add(c.duration)) {
		c.Fetch(ctx)
	}
	return c.cache
}
