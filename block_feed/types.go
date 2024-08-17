package block_feed

import (
	cometbft "github.com/cometbft/cometbft/types"
)

// BlockFeed is a standard interface to provide subscription over blocks
// There is only one method OnBlockFound and it gives you access to the
// BlockFeed channel
type BlockFeed interface {
	// Close closes underlying subscriber
	Close() error

	// Subscribe starts subscription to the block source
	Subscribe(rpcIndex int) (chan *BlockResult, error)
}

type BlockResult struct {
	BlockID *cometbft.BlockID `json:"block_id"`
	Block   *cometbft.Block   `json:"block"`
}
