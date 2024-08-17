package tx

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/osmosis-labs/mantlemint/db/safe_batch"
	"github.com/osmosis-labs/mantlemint/mantlemint"
	"github.com/stretchr/testify/assert"

	cbftdb "github.com/cometbft/cometbft-db"
	cbftjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/types"
)

func TestIndexTx(t *testing.T) {
	db := cbftdb.NewMemDB()
	block := &types.Block{}
	blockFile, _ := os.Open("../fixtures/block_4814775.json")
	blockJSON, _ := ioutil.ReadAll(blockFile)
	if err := cbftjson.Unmarshal(blockJSON, block); err != nil {
		t.Fail()
	}

	eventFile, _ := os.Open("../fixtures/response_4814775.json")
	eventJSON, _ := io.ReadAll(eventFile)
	evc := mantlemint.NewMantlemintEventCollector()
	event := types.EventDataTx{}
	if err := cbftjson.Unmarshal(eventJSON, &event.Result); err != nil {
		panic(err)
	}

	_ = evc.PublishEventTx(event)

	safebatch := safe_batch.NewSafeBatchDB(db)
	if err := IndexTx(*safebatch.(*safe_batch.SafeBatchDB), block, nil, evc, nil); err != nil {
		panic(err)
	}
	safebatch.(safe_batch.SafeBatchDBCloser).Flush()

	txn, err := txByHashHandler(db, "C794D5CE7179AED455C10E8E7645FE8F8A40BA0C97F1275AB87B5E88A52CB2C3")
	assert.Nil(t, err)
	assert.NotNil(t, txn)
	fmt.Println(string(txn))

	txns, err := txsByHeightHandler(db, "4814775")
	assert.Nil(t, err)
	assert.NotNil(t, txns)
	fmt.Println(string(txns))
}
