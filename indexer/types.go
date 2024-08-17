package indexer

import (
	"log"
	"net/http"
	"runtime"

	"github.com/gorilla/mux"

	cbftdb "github.com/cometbft/cometbft-db"
	cbfttypes "github.com/cometbft/cometbft/types"
	"github.com/osmosis-labs/mantlemint/db/safe_batch"
	"github.com/osmosis-labs/mantlemint/mantlemint"
	"github.com/osmosis-labs/osmosis/v25/app"
)

type IndexFunc func(indexerDB safe_batch.SafeBatchDB, block *cbfttypes.Block, blockId *cbfttypes.BlockID, evc *mantlemint.EventCollector, app *app.OsmosisApp) error
type ClientHandler func(w http.ResponseWriter, r *http.Request) error
type RESTRouteRegisterer func(router *mux.Router, indexerDB cbftdb.DB)

func CreateIndexer(idf IndexFunc) IndexFunc {
	return idf
}

func CreateRESTRoute(registerer RESTRouteRegisterer) RESTRouteRegisterer {
	return registerer
}

var (
	ErrorInternal = func(err error) string {
		_, fn, fl, ok := runtime.Caller(1)

		if !ok {
			// ...
		} else {
			log.Printf("ErrorInternal[%s:%d] %v\n", fn, fl, err.Error())
		}

		return "internal server error"
	}
)
