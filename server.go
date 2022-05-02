package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
    "database/sql"

	"github.com/gorilla/mux"
	"github.com/temelpa/timetravel/api"
	"github.com/temelpa/timetravel/service"
    _ "github.com/mattn/go-sqlite3"
)

// logError logs all non-nil errors
func logError(err error) {
	if err != nil {
		log.Printf("error: %v", err)
	}
}

func initDb()(*sql.DB, error) {
    db, err := sql.Open("sqlite3", "./records.db")
    if err != nil {
        logError(err)
        return nil, err
    }

    create := `
        CREATE TABLE IF NOT EXISTS records (
            id INTEGER NOT NULL PRIMARY KEY,
            rid INTEGER NOT NULL,
            key TEXT,
            value TEXT);
    `
    _, err = db.Exec(create);

    if err != nil {
        logError(err)
        return nil, err
    }
    // Check to see if index exists.
    indexSearch := `
        PRAGMA index_list(records);
    `
    rows, err := db.Query(indexSearch);

    // If index does not exist, add one to index recordId and id to
    // improve the performance of the query
    if !rows.Next() {
        index := `
            CREATE INDEX recordId
            ON records (rid,id);
        `
        _, err = db.Exec(index);

        if err != nil {
            logError(err)
            return nil, err
        }
    }

    return db, nil
}


func main() {
	router := mux.NewRouter()

    db, err := initDb()
    if err != nil {
        logError(err)
        return
    }

	service := service.NewInMemoryRecordService()
	apiv1 := api.NewAPI(&service)
	apiRoute := router.PathPrefix("/api/v1").Subrouter()
	apiRoute.Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		logError(err)
	})
	apiv1.CreateRoutesV1(apiRoute)

	apiv2 := api.NewAPIv2(db)
	apiRoutev2 := router.PathPrefix("/api/v2").Subrouter()
	apiRoutev2.Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		logError(err)
	})
	apiv2.CreateRoutesV2(apiRoutev2)

	address := "127.0.0.1:8000"
	srv := &http.Server{
		Handler:      router,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("listening on %s", address)
	log.Fatal(srv.ListenAndServe())
    defer db.Close()
}
