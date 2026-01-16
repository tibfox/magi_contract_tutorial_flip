package main

import (
	"magi_contract_tutorial_flip/sdk"
	"strconv"
	"strings"
)

func main() {
	// needs to be here for wasm export to work
}

// Flip randomly shuffles a list of possibilities provided in the payload
// Payload format: "possibility1|possibility2|possibility3|..."
// Returns: "storageKey -> possibility3|possibility1|..."
// Example: "hive:flipuser:12345678-0-0 -> possibility3|possibility1|..."
// Logs: "flip;sender:<user>;result:possibility3|possibility1|..."
// Stores result under key: "sender:blockHeight-index-opIndex"
// Example key: "hive:flipuser:12345678-0-0"

//go:wasmexport flip
func Flip(payload *string) *string {
	// Check if payload is provided
	if payload == nil || *payload == "" {
		sdk.Abort("payload is required")
		return nil
	}

	// Parse the |-delimited payload
	possibilities := strings.Split(*payload, "|")
	if len(possibilities) < 2 {
		sdk.Abort("at least 2 possibilities are required")
		return nil
	}

	// Get environment info
	env := sdk.GetEnv()
	sender := env.Sender.Address.String()

	// Generate a seed from call data for randomization
	// Using BlockHeight-Index-OpIndex ensures uniqueness for each operation
	// but deterministic for the same operation
	seed := generateSeed(env.BlockHeight, env.Index, env.OpIndex)

	// Shuffle the possibilities using Fisher-Yates algorithm
	shuffled := shuffleWithSeed(possibilities, seed)

	// Join the shuffled result
	result := strings.Join(shuffled, "|")

	// Create storage key: blockHeight-index-opIndex
	storageKey := strconv.FormatUint(env.BlockHeight, 10) +
		"-" + strconv.FormatUint(env.Index, 10) +
		"-" + strconv.FormatUint(env.OpIndex, 10)

	// Store the result
	sdk.StateSetObject(storageKey, result)

	// Log the result: flip;sender:<user>;result:possiblities|list|of
	logMessage := "flip;sender:" + sender + ";result:" + result
	sdk.Log(logMessage)

	// Return "storageKey -> possiblities|list|of"
	returnValue := storageKey + " -> " + result
	return &returnValue
}
