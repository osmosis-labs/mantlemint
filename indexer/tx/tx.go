package tx

import (
	"encoding/json"
	"fmt"

	cbftjson "github.com/cometbft/cometbft/libs/json"
	cbfttypes "github.com/cometbft/cometbft/types"
	"github.com/osmosis-labs/mantlemint/db/safe_batch"
	"github.com/osmosis-labs/mantlemint/indexer"
	"github.com/osmosis-labs/mantlemint/mantlemint"
	"github.com/osmosis-labs/osmosis/v25/app"
)

var cdc = app.MakeEncodingConfig()

var IndexTx = indexer.CreateIndexer(func(batch safe_batch.SafeBatchDB, block *cbfttypes.Block, blockID *cbfttypes.BlockID, evc *mantlemint.EventCollector, _ *app.OsmosisApp) error {
	// encoder; proto -> mem -> json
	txDecoder := cdc.TxConfig.TxDecoder()
	jsonEncoder := cdc.TxConfig.TxJSONEncoder()

	txHashes := make([]string, len(block.Txs))
	txRecords := make([]TxRecord, len(block.Txs))
	byHeightPayload := make([]TxByHeightRecord, len(block.Txs))

	// by hash
	for txIndex, txByte := range block.Txs {
		txRecord := TxRecord{}

		hash := txByte.Hash()
		tx, decodeErr := txDecoder(txByte)

		if decodeErr != nil {
			return decodeErr
		}

		// encode tx to JSON for max compat & shave deserialization cost at serving
		txJSON, _ := jsonEncoder(tx)

		// handle response -> json
		response := ToResponseDeliverTxJSON(evc.ResponseDeliverTxs[txIndex])
		responseJSON, responseMarshalErr := cbftjson.Marshal(response)

		if responseMarshalErr != nil {
			return responseMarshalErr
		}

		// populate txRecord
		txRecord.Tx = txJSON
		txRecord.TxResponse = responseJSON

		txHashes[txIndex] = fmt.Sprintf("%X", hash)
		txRecords[txIndex] = txRecord

		// byHeightRecord
		// handle non-successful case first
		byHeightPayload[txIndex].Code = response.Code
		byHeightPayload[txIndex].Codespace = response.Codespace
		byHeightPayload[txIndex].GasUsed = response.GasUsed
		byHeightPayload[txIndex].GasWanted = response.GasWanted
		byHeightPayload[txIndex].Height = block.Height
		byHeightPayload[txIndex].RawLog = response.Log
		byHeightPayload[txIndex].Logs = func() json.RawMessage {
			if response.Code == 0 {
				return []byte(response.Log)
			} else {
				out, _ := json.Marshal([]string{})
				return out
			}
		}()
		byHeightPayload[txIndex].TxHash = fmt.Sprintf("%X", hash)
		byHeightPayload[txIndex].Timestamp = block.Time
		byHeightPayload[txIndex].Tx = txJSON
	}

	// 1. byHash -- matching the interface for /cosmos/tx/v1beta1/txs/{hash}
	for txIndex, txRecord := range txRecords {
		txRecordJSON, marshalErr := cbftjson.Marshal(txRecord)
		if marshalErr != nil {
			return marshalErr
		}

		batchSetErr := batch.Set(getKey(txHashes[txIndex]), txRecordJSON)
		if batchSetErr != nil {
			return batchSetErr
		}
	}

	// 2. byHeight -- custom endpoint
	byHeightJSON, byHeightErr := cbftjson.Marshal(byHeightPayload)
	if byHeightErr != nil {
		return byHeightErr
	}

	batchSetErr := batch.Set(getByHeightKey(uint64(block.Height)), byHeightJSON)
	if batchSetErr != nil {
		return batchSetErr
	}

	return nil
})
