package utils

import (
	"errors"
	"fmt"

	"github.com/gogo/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"

	db "github.com/cometbft/cometbft-db"
)

func GatherAllKeysFromStore(storeObj store.KVStore) []string {
	iterator := storeObj.Iterator(nil, nil)
	defer iterator.Close()

	keys := []string{}
	for ; iterator.Valid(); iterator.Next() {
		keys = append(keys, string(iterator.Key()))
	}
	return keys
}

func GatherValuesFromStore[T any](storeObj store.KVStore, keyStart, keyEnd []byte, parseValue func([]byte) (T, error)) ([]T, error) {
	iterator := storeObj.Iterator(keyStart, keyEnd)
	defer iterator.Close()
	return gatherValuesFromIteratorWithStop(iterator, parseValue, noStopFn)
}

func GatherValuesFromStorePrefix[T any](storeObj store.KVStore, prefix []byte, parseValue func([]byte) (T, error)) ([]T, error) {
	iterator := sdk.KVStorePrefixIterator(storeObj, prefix)
	defer iterator.Close()
	return gatherValuesFromIteratorWithStop(iterator, parseValue, noStopFn)
}

func GetValuesUntilDerivedStop[T any](storeObj store.KVStore, keyStart []byte, stopFn func([]byte) bool, parseValue func([]byte) (T, error)) ([]T, error) {
	// SDK iterator is broken for nil end time, and non-nil start time
	// https://github.com/cosmos/cosmos-sdk/issues/12661
	// hence we use []byte{0xff}
	keyEnd := []byte{0xff}
	return GetIterValuesWithStop(storeObj, keyStart, keyEnd, false, stopFn, parseValue)
}

func makeIterator(storeObj store.KVStore, keyStart, keyEnd []byte, reverse bool) store.Iterator {
	if reverse {
		return storeObj.ReverseIterator(keyStart, keyEnd)
	}
	return storeObj.Iterator(keyStart, keyEnd)
}

func GetIterValuesWithStop[T any](
	storeObj store.KVStore,
	keyStart []byte,
	keyEnd []byte,
	reverse bool,
	stopFn func([]byte) bool,
	parseValue func([]byte) (T, error),
) ([]T, error) {
	iter := makeIterator(storeObj, keyStart, keyEnd, reverse)
	defer iter.Close()

	return gatherValuesFromIteratorWithStop(iter, parseValue, stopFn)
}

func GetFirstValueAfterPrefixInclusive[T any](storeObj store.KVStore, keyStart []byte, parseValue func([]byte) (T, error)) (T, error) {
	// SDK iterator is broken for nil end time, and non-nil start time
	// https://github.com/cosmos/cosmos-sdk/issues/12661
	// hence we use []byte{0xff}
	return GetFirstValueInRange(storeObj, keyStart, []byte{0xff}, false, parseValue)
}

func GetFirstValueInRange[T any](storeObj store.KVStore, keyStart, keyEnd []byte, reverseIterate bool, parseValue func([]byte) (T, error)) (T, error) {
	iterator := makeIterator(storeObj, keyStart, keyEnd, reverseIterate)
	defer iterator.Close()

	if !iterator.Valid() {
		var blankValue T
		return blankValue, errors.New("no values in range")
	}

	return parseValue(iterator.Value())
}

func gatherValuesFromIteratorWithStop[T any](iterator db.Iterator, parseValue func([]byte) (T, error), stopFn func([]byte) bool) ([]T, error) {
	values := []T{}
	for ; iterator.Valid(); iterator.Next() {
		if stopFn(iterator.Key()) {
			break
		}
		val, err := parseValue(iterator.Value())
		if err != nil {
			return nil, err
		}
		values = append(values, val)
	}
	return values, nil
}

func noStopFn([]byte) bool {
	return false
}

// MustSet runs store.Set(key, proto.Marshal(value))
// but panics on any error.
func MustSet(storeObj store.KVStore, key []byte, value proto.Message) {
	bz, err := proto.Marshal(value)
	if err != nil {
		panic(err)
	}

	storeObj.Set(key, bz)
}

// MustGet gets key from store by mutating result
// Panics on any error.
func MustGet(storeObj store.KVStore, key []byte, result proto.Message) {
	b := storeObj.Get(key)
	if b == nil {
		panic(fmt.Errorf("getting at key (%v) should not have been nil", key))
	}
	if err := proto.Unmarshal(b, result); err != nil {
		panic(err)
	}
}

// MustSetDec sets dec value to store at key. Panics on any error.
func MustSetDec(storeObj store.KVStore, key []byte, value sdk.Dec) {
	MustSet(storeObj, key, &sdk.DecProto{
		Dec: value,
	})
}

// MustGetDec gets dec value from store at key. Panics on any error.
func MustGetDec(storeObj store.KVStore, key []byte) sdk.Dec {
	result := &sdk.DecProto{}
	MustGet(storeObj, key, result)
	return result.Dec
}
