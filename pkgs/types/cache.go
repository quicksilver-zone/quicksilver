package types

import (
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
	Init(url string, updateTime time.Duration)
	Fetch()
	Get() []T
}

var _ CacheManagerI[any] = &CacheManager[any]{}

type CacheManager[T any] struct {
	url         string
	cache       []T
	lastUpdated time.Time
	duration    time.Duration
}

func (c *CacheManager[T]) Init(url string, updateInterval time.Duration) {
	c.url = url
	c.duration = updateInterval
	c.Fetch()
}

func (c *CacheManager[T]) Fetch() {
	fmt.Println("Fetching and caching " + c.url)
	response, err := http.Get(c.url)
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

func (c CacheManager[T]) Get() []T {
	if time.Now().After(c.lastUpdated.Add(c.duration)) {
		c.Fetch()
	}
	return c.cache
}
