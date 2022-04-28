package api

import (
    "database/sql"

	"github.com/gorilla/mux"
	"github.com/temelpa/timetravel/service"
    _ "github.com/mattn/go-sqlite3"
)

type API struct {
	records service.RecordService
    db *sql.DB
}

func NewAPI(records service.RecordService, db *sql.DB) *API {
	return &API{
            records: records,
            db: db,
    }
}

// generates all api routes
func (a *API) CreateRoutes(routes *mux.Router) {
	routes.Path("/records/{id}").HandlerFunc(a.GetRecords).Methods("GET")
	routes.Path("/records/{id}").HandlerFunc(a.PostRecords).Methods("POST")
}
