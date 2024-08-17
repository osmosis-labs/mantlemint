package block

import (
	cbfttypes "github.com/cometbft/cometbft/types"
	"github.com/osmosis-labs/mantlemint/lib"
)

var prefix = []byte("block/height:")
var getKey = func(height uint64) []byte {
	return lib.ConcatBytes(prefix, lib.UintToBigEndian(height))
}

type BlockRecord struct {
	BlockID *cbfttypes.BlockID `json:"block_id"`
	Block   *cbfttypes.Block   `json:"block"`
}
