package api

import (
    "database/sql"

	"github.com/gorilla/mux"
    "github.com/temelpa/timetravel/service"
    _ "github.com/mattn/go-sqlite3"
)

type API struct {
    records service.RecordService
}

func NewAPI(records service.RecordService) *API {
	return &API{records}
}

// generates all api routes for v1 route
func (a *API) CreateRoutesV1(routes *mux.Router) {
	routes.Path("/records/{id}").HandlerFunc(a.GetRecords).Methods("GET")
	routes.Path("/records/{id}").HandlerFunc(a.PostRecords).Methods("POST")
}

type APIv2 struct {
    eid int
    db *sql.DB
}

func NewAPIv2(db *sql.DB) *APIv2 {
	return &APIv2{
            eid: 0,
            db: db,
    }
}

func (a *APIv2) IncrementEid() {
    a.eid++
}

// generates all api routes for v2 route
func (a *APIv2) CreateRoutesV2(routes *mux.Router) {
	routes.Path("/records/{id}").HandlerFunc(a.GetRecordsV2).Methods("GET")
	routes.Path("/records/{id}").HandlerFunc(a.PostRecordsV2).Methods("POST")

    // new route to look at history of an rid, key pair
	routes.Path("/records/{id}/{key}").HandlerFunc(a.GetRecordForKeyV2).Methods("GET")
}

