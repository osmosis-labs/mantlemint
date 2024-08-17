package rootmulti

import (
	"io"

	"github.com/cosmos/cosmos-sdk/store/dbadapter"
	pruningtypes "github.com/cosmos/cosmos-sdk/store/pruning/types"
	"github.com/cosmos/cosmos-sdk/store/types"

	dbm "github.com/cometbft/cometbft-db"
)

var commithash = []byte("FAKE_HASH")

//----------------------------------------
// commitDBStoreWrapper should only be used for simulation/debugging,
// as it doesn't compute any commit hash, and it cannot load older state.

// Wrapper type for dbm.Db with implementation of KVStore
type commitDBStoreAdapter struct {
	dbadapter.Store
	prefix []byte
}

// CacheWrap implements types.CommitKVStore.
// Subtle: this method shadows the method (Store).CacheWrap of commitDBStoreAdapter.Store.
func (cdsa commitDBStoreAdapter) CacheWrap() types.CacheWrap {
	return cdsa.Store.CacheWrap()
}

// CacheWrapWithTrace implements types.CommitKVStore.
// Subtle: this method shadows the method (Store).CacheWrapWithTrace of commitDBStoreAdapter.Store.
func (cdsa commitDBStoreAdapter) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	return cdsa.Store.CacheWrapWithTrace(w, tc)
}

// Delete implements types.CommitKVStore.
// Subtle: this method shadows the method (Store).Delete of commitDBStoreAdapter.Store.
func (cdsa commitDBStoreAdapter) Delete(key []byte) {
	if err := cdsa.DB.Delete(key); err != nil {
		panic(err)
	}
}

// Get implements types.CommitKVStore.
// Subtle: this method shadows the method (Store).Get of commitDBStoreAdapter.Store.
func (cdsa commitDBStoreAdapter) Get(key []byte) []byte {
	value, err := cdsa.DB.Get(key)
	if err != nil {
		panic(err)
	}
	return value
}

// GetStoreType implements types.CommitKVStore.
// Subtle: this method shadows the method (Store).GetStoreType of commitDBStoreAdapter.Store.
func (cdsa commitDBStoreAdapter) GetStoreType() types.StoreType {
	return types.StoreTypeDB
}

// Has implements types.CommitKVStore.
// Subtle: this method shadows the method (Store).Has of commitDBStoreAdapter.Store.
func (cdsa commitDBStoreAdapter) Has(key []byte) bool {
	has, err := cdsa.DB.Has(key)
	if err != nil {
		panic(err)
	}
	return has
}

// Iterator implements types.CommitKVStore.
// Subtle: this method shadows the method (Store).Iterator of commitDBStoreAdapter.Store.
func (cdsa commitDBStoreAdapter) Iterator(start []byte, end []byte) dbm.Iterator {
	itr, err := cdsa.DB.Iterator(start, end)
	if err != nil {
		panic(err)
	}
	return itr
}

// ReverseIterator implements types.CommitKVStore.
// Subtle: this method shadows the method (Store).ReverseIterator of commitDBStoreAdapter.Store.
func (cdsa commitDBStoreAdapter) ReverseIterator(start []byte, end []byte) dbm.Iterator {
	itr, err := cdsa.DB.ReverseIterator(start, end)
	if err != nil {
		panic(err)
	}
	return itr
}

// Set implements types.CommitKVStore.
// Subtle: this method shadows the method (Store).Set of commitDBStoreAdapter.Store.
func (cdsa commitDBStoreAdapter) Set(key []byte, value []byte) {
	cdsa.DB.Set(key, value)
}

// SetCommitting implements types.CommitKVStore.
func (cdsa commitDBStoreAdapter) SetCommitting() {
}

// UnsetCommitting implements types.CommitKVStore.
func (cdsa commitDBStoreAdapter) UnsetCommitting() {
}

func (cdsa commitDBStoreAdapter) Commit() types.CommitID {
	return types.CommitID{
		Version: -1,
		Hash:    commithash,
	}
}

func (cdsa commitDBStoreAdapter) LastCommitID() types.CommitID {
	return types.CommitID{
		Version: -1,
		Hash:    commithash,
	}
}

func (cdsa commitDBStoreAdapter) SetPruning(_ pruningtypes.PruningOptions) {}

// GetPruning is a no-op as pruning options cannot be directly set on this store.
// They must be set on the root commit multi-store.
func (cdsa commitDBStoreAdapter) GetPruning() pruningtypes.PruningOptions {
	return pruningtypes.PruningOptions{}
}

func (cdsa *commitDBStoreAdapter) BranchStoreWithHeightLimitedDB(hldb dbm.DB) types.CommitKVStore {
	var db = dbm.NewPrefixDB(hldb, cdsa.prefix)

	return commitDBStoreAdapter{Store: dbadapter.Store{DB: db}, prefix: cdsa.prefix}
}
