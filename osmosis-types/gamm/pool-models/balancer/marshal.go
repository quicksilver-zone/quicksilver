package balancer

import (
	"encoding/json"
)

// Upstream marshal.go contains balancerPoolPretty and custom
// JSONMarshal and JSONUnmarshal functions. Given that we deal
// directly with the value in the KV store, not the output from
// a query, we do not need to use these methods and they are
// removed below.

func (p Pool) String() string {
	out, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(out)
}
