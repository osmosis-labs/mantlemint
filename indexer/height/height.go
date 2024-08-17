package height

import (
	"fmt"

	cbftjson "github.com/cometbft/cometbft/libs/json"
	cbfttypes "github.com/cometbft/cometbft/types"
	"github.com/osmosis-labs/mantlemint/db/safe_batch"
	"github.com/osmosis-labs/mantlemint/indexer"
	"github.com/osmosis-labs/mantlemint/mantlemint"
	"github.com/osmosis-labs/osmosis/v25/app"
)

var IndexHeight = indexer.CreateIndexer(func(indexerDB safe_batch.SafeBatchDB, block *cbfttypes.Block, _ *cbfttypes.BlockID, _ *mantlemint.EventCollector, _ *app.OsmosisApp) error {
	defer fmt.Printf("[indexer/height] indexing done for height %d\n", block.Height)
	height := block.Height

	record := HeightRecord{Height: uint64(height)}
	recordJSON, recordErr := cbftjson.Marshal(record)
	if recordErr != nil {
		return recordErr
	}

	return indexerDB.Set(getKey(), recordJSON)
})
