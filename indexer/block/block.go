package block

import (
	"fmt"

	// tmjson "github.com/tendermint/tendermint/libs/json"
	// tm "github.com/tendermint/tendermint/types"
	// terra "github.com/terra-money/core/v2/app"
	// "github.com/terra-money/mantlemint/db/safe_batch"
	// "github.com/terra-money/mantlemint/indexer"
	// "github.com/terra-money/mantlemint/mantlemint"
)

var IndexBlock = indexer.CreateIndexer(func(indexerDB safe_batch.SafeBatchDB, block *tm.Block, blockID *tm.BlockID, _ *mantlemint.EventCollector, _ *terra.TerraApp) error {
	defer fmt.Printf("[indexer/block] indexing done for height %d\n", block.Height)
	record := BlockRecord{
		Block:   block,
		BlockID: blockID,
	}

	recordJSON, recordErr := tmjson.Marshal(record)
	if recordErr != nil {
		return recordErr
	}

	return indexerDB.Set(getKey(uint64(block.Height)), recordJSON)
})
