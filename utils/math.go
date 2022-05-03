package utils

import "math"

func MinI64(in []int64) int64 {
	out := int64(math.MaxInt64)
	for _, i := range in {
		if i < out {
			out = i
		}
	}
	return out
}

func MaxI64(in []int64) int64 {
	out := int64(math.MinInt64)
	for _, i := range in {
		if i > out {
			out = i
		}
	}
	return out
}

func MinU64(in []uint64) uint64 {
	out := uint64(math.MaxUint64)
	for _, i := range in {
		if i < out {
			out = i
		}
	}
	return out
}

func MaxU64(in []uint64) uint64 {
	out := uint64(0)
	for _, i := range in {
		if i > out {
			out = i
		}
	}
	return out
}
