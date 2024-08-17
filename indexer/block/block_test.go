package block

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/osmosis-labs/mantlemint/db/safe_batch"
	"github.com/stretchr/testify/assert"

	cbftdb "github.com/cometbft/cometbft-db"
	cbftjson "github.com/cometbft/cometbft/libs/json"
)

func TestIndexBlock(t *testing.T) {
	db := cbftdb.NewMemDB()
	blockFile, _ := os.Open("../fixtures/block_4724005_raw.json")
	blockJSON, _ := io.ReadAll(blockFile)

	record := BlockRecord{}
	_ = cbftjson.Unmarshal(blockJSON, &record)

	batch := safe_batch.NewSafeBatchDB(db)
	batch.(safe_batch.SafeBatchDBCloser).Open()
	if err := IndexBlock(*batch.(*safe_batch.SafeBatchDB), record.Block, record.BlockID, nil, nil); err != nil {
		panic(err)
	}
	batch.(safe_batch.SafeBatchDBCloser).Flush()

	block, err := blockByHeightHandler(db, "4724005")
	assert.Nil(t, err)
	assert.NotNil(t, block)

	fmt.Println(string(block))
}
