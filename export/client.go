package export

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/osmosis-labs/osmosis/v25/app"
)

func RegisterRESTRoutes(router *mux.Router, app *app.OsmosisApp) {
	router.Handle("/export/accounts", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		err := ExportAllAccounts(app)
		if err != nil {
			writer.WriteHeader(http.StatusConflict)
			writer.Write([]byte(err.Error()))
		}
		writer.WriteHeader(http.StatusOK)
	})).Methods("POST")

	router.Handle("/export/circulating_supply", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cs, err := ExportCirculatingSupply(app)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(err.Error()))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cs.String()))
	}))
}
