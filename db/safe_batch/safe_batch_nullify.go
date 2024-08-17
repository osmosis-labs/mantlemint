package safe_batch

import cbftdb "github.com/cometbft/cometbft-db"

var _ cbftdb.Batch = (*SafeBatchNullified)(nil)

type SafeBatchNullified struct {
	batch cbftdb.Batch
}

func NewSafeBatchNullify(batch cbftdb.Batch) cbftdb.Batch {
	return &SafeBatchNullified{
		batch: batch,
	}
}

func (s SafeBatchNullified) Set(key, value []byte) error {
	return s.batch.Set(key, value)
}

func (s SafeBatchNullified) Delete(key []byte) error {
	return s.batch.Delete(key)
}

func (s SafeBatchNullified) Write() error {
	// noop
	return nil
}

func (s SafeBatchNullified) WriteSync() error {
	return s.Write()
}

func (s SafeBatchNullified) Close() error {
	// noop
	return nil
}
