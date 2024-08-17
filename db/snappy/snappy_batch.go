package snappy

import (
	cometbft "github.com/cometbft/cometbft-db"
	"github.com/golang/snappy"
)

var _ cometbft.Batch = (*SnappyBatch)(nil)

type SnappyBatch struct {
	batch cometbft.Batch
}

func NewSnappyBatch(batch cometbft.Batch) *SnappyBatch {
	return &SnappyBatch{
		batch: batch,
	}
}

func (s *SnappyBatch) Set(key, value []byte) error {
	return s.batch.Set(key, snappy.Encode(nil, value))
}

func (s *SnappyBatch) Delete(key []byte) error {
	return s.batch.Delete(key)
}

func (s *SnappyBatch) Write() error {
	return s.batch.Write()
}

func (s *SnappyBatch) WriteSync() error {
	return s.batch.WriteSync()
}

func (s *SnappyBatch) Close() error {
	return s.batch.Close()
}
