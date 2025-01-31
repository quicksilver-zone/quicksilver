package main

import (
	"fmt"
	"maps"
	"math/rand"
)

// DistributeConfigs ensures an even spread of configurations across nodes
func DistributeConfigs(N, M, redundancy, seed int) map[int][]int {
	nodeConfigs := make(map[int][]int)

	// Step 1: Create a list with `redundancy` copies of each config
	configPool := make([]int, 0, N*redundancy)
	for i := 0; i < N; i++ {
		for j := 0; j < redundancy; j++ {
			configPool = append(configPool, i)
		}
	}

	// Step 2: Shuffle the config pool with a fixed seed for deterministic results
	rand.Seed(int64(seed))
	rand.Shuffle(len(configPool), func(i, j int) { configPool[i], configPool[j] = configPool[j], configPool[i] })

	// Step 3: Assign configurations to nodes in a round-robin fashion
	for i, config := range configPool {
		node := i % M // Ensures an even spread
		nodeConfigs[node] = append(nodeConfigs[node], config)
	}

	return nodeConfigs
}

func CheckDistributions(dist map[int][]int) bool {
	for i := range maps.Values(dist) {
		if !IsUniqueSliceElements(i) { return false; }
	}
	return true
}

func IsUniqueSliceElements[T comparable](inputSlice []T) bool {
	seen := make(map[T]bool, len(inputSlice))
	for _, element := range inputSlice {
		if seen[element] {
			return false
		}
		seen[element] = true
	}
	return true
}

func main() {
	N := 10  // Number of configurations
	M := 4   // Number of nodes
	redundancy := 3 // Each config appears exactly `redundancy` times
	seed := 0 // Deterministic seed for reproducibility

	var result map[int][]int
	for true {
		seed = seed + 1
		result = DistributeConfigs(N, M, redundancy, seed)
		check := CheckDistributions(result)
		if check { break }
	}

	fmt.Printf("Took %d iterations\n\n", seed) 
	for node, configs := range result {
		fmt.Printf("Node %d will have configs: %v (Total: %d)\n", node, configs, len(configs))
	}
}

