package test

import (
	"encoding/json"
	"os"
	"testing"

	"vsc-node/lib/test_utils"
	"vsc-node/modules/db/vsc/contracts"
	ledgerDb "vsc-node/modules/db/vsc/ledger"
	stateEngine "vsc-node/modules/state-processing"

	"github.com/stretchr/testify/assert"
)

func TestFlip(t *testing.T) {
	testCode, err := os.ReadFile("../artifacts/main.wasm")
	assert.NoError(t, err)

	ct := test_utils.NewContractTest()
	ct.Deposit("hive:flipuser", 5000, ledgerDb.AssetHive)
	ct.RegisterContract("vsc_flip_contract", "hive:flipowner", testCode)

	txSelf := stateEngine.TxSelf{
		TxId:                 "2testtxid",
		BlockId:              "blockid122",
		Index:                0,
		OpIndex:              0,
		Timestamp:            "2025-09-03T12:00:00",
		RequiredAuths:        []string{"hive:flipuser"},
		RequiredPostingAuths: []string{},
	}

	callResult, _, logsMap := ct.Call(stateEngine.TxVscCallContract{
		Self:       txSelf,
		ContractId: "vsc_flip_contract",
		Action:     "flip",
		Payload:    json.RawMessage([]byte(`"tibfox|stevenson7|meno|ph1102"`)),
		RcLimit:    1000,
		Intents:    []contracts.Intent{},
	})

	t.Logf("Return: %s", callResult.Ret)
	t.Logf("RC: %.3f", float64(callResult.RcUsed)/1000)
	t.Logf("Logs: %v", logsMap["vsc_flip_contract"])

	assert.True(t, callResult.Success, "flip call should succeed")
}
