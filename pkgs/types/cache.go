package types

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Data[T any] struct {
	Data []T
}

type CacheManagerI[T any] interface {
	Init(ctx context.Context, url string, updateTime time.Duration)
	Fetch(ctx context.Context)
	Get(ctx context.Context) []T
}

var _ CacheManagerI[any] = &CacheManager[any]{}

type CacheManager[T any] struct {
	url         string
	cache       []T
	lastUpdated time.Time
	duration    time.Duration
}

func (c *CacheManager[T]) Init(ctx context.Context, url string, updateInterval time.Duration) {
	c.url = url
	c.duration = updateInterval
	c.Fetch(ctx)
}

func (c *CacheManager[T]) Fetch(ctx context.Context) {
	fmt.Println("Fetching and caching " + c.url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, http.NoBody)
	if err != nil {
		panic(err)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	data := Data[T]{}

	err = json.Unmarshal(responseData, &data)
	if err != nil {
		panic(err)
	}
	c.cache = data.Data
	c.lastUpdated = time.Now()
}

func (c CacheManager[T]) Get(ctx context.Context) []T {
	if time.Now().After(c.lastUpdated.Add(c.duration)) {
		c.Fetch(ctx)
	}
	return c.cache
}
