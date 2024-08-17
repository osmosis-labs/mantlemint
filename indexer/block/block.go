package block

import (
	"fmt"

	cbftjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/types"
	"github.com/osmosis-labs/mantlemint/db/safe_batch"
	"github.com/osmosis-labs/mantlemint/indexer"
	"github.com/osmosis-labs/mantlemint/mantlemint"
	"github.com/osmosis-labs/osmosis/v25/app"
)

var IndexBlock = indexer.CreateIndexer(func(indexerDB safe_batch.SafeBatchDB, block *types.Block, blockID *types.BlockID, _ *mantlemint.EventCollector, _ *app.OsmosisApp) error {
	defer fmt.Printf("[indexer/block] indexing done for height %d\n", block.Height)
	record := BlockRecord{
		Block:   block,
		BlockID: blockID,
	}

	recordJSON, recordErr := cbftjson.Marshal(record)
	if recordErr != nil {
		return recordErr
	}

	return indexerDB.Set(getKey(uint64(block.Height)), recordJSON)
})
