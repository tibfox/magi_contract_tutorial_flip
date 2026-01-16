package main

// generateSeed creates a numeric seed from block height and indices
// BlockHeight-Index-OpIndex combination is unique for each operation
func generateSeed(blockHeight uint64, index uint64, opIndex uint64) uint64 {
	// Combine the three values into a single seed
	// Using prime multipliers to reduce collision likelihood
	return blockHeight*1000003 + index*1009 + opIndex
}

// shuffleWithSeed performs Fisher-Yates shuffle with a deterministic seed
func shuffleWithSeed(items []string, seed uint64) []string {
	n := len(items)
	result := make([]string, n)
	copy(result, items)

	// Simple LCG random number generator
	// Using parameters from Numerical Recipes
	a := uint64(1664525)
	c := uint64(1013904223)
	m := uint64(1 << 32)
	state := seed

	for i := n - 1; i > 0; i-- {
		// Generate next random number
		state = (a*state + c) % m
		j := int(state % uint64(i+1))

		// Swap
		result[i], result[j] = result[j], result[i]
	}

	return result
}
