package rollbackable

import (
	cbftdb "github.com/cometbft/cometbft-db"
)

type HasRollbackBatch interface {
	RollbackBatch() cbftdb.Batch
}

var _ cbftdb.Batch = (*RollbackableBatch)(nil)

type RollbackableBatch struct {
	cbftdb.Batch

	db            cbftdb.DB
	RollbackBatch cbftdb.Batch
	RecordCount   int
}

func NewRollbackableBatch(db cbftdb.DB) *RollbackableBatch {
	return &RollbackableBatch{
		db:            db,
		Batch:         db.NewBatch(),
		RollbackBatch: db.NewBatch(),
	}
}

// revert value for key to previous state
func (b *RollbackableBatch) backup(key []byte) error {
	b.RecordCount++
	data, err := b.db.Get(key)
	if err != nil {
		return err
	}
	if data == nil {
		return b.RollbackBatch.Delete(key)
	} else {
		return b.RollbackBatch.Set(key, data)
	}
}

func (b *RollbackableBatch) Set(key, value []byte) error {
	if err := b.backup(key); err != nil {
		return err
	}
	return b.Batch.Set(key, value)
}

func (b *RollbackableBatch) Delete(key []byte) error {
	if err := b.backup(key); err != nil {
		return err
	}
	return b.Batch.Delete(key)
}
