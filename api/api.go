package api

import (
    "database/sql"

	"github.com/gorilla/mux"
    _ "github.com/mattn/go-sqlite3"
)

type API struct {
    eid int
    db *sql.DB
}

func NewAPI(db *sql.DB) *API {
	return &API{
            eid: 0,
            db: db,
    }
}

func (a *API) IncrementEid() {
    a.eid++
}

// generates all api routes
func (a *API) CreateRoutes(routes *mux.Router) {
	routes.Path("/records/{id}").HandlerFunc(a.GetRecords).Methods("GET")
	routes.Path("/records/{id}").HandlerFunc(a.PostRecords).Methods("POST")
}
