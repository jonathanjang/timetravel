package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
    "database/sql"

	"github.com/gorilla/mux"
	"github.com/temelpa/timetravel/api"
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

    // TODO: add indexing via rid
    const create string = `
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

    return db, nil
}


func main() {
	router := mux.NewRouter()

    db, err := initDb()
    if err != nil {
        logError(err)
        return
    }

	api := api.NewAPI(db)
	apiRoute := router.PathPrefix("/api/v1").Subrouter()
	apiRoute.Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		logError(err)
	})
	api.CreateRoutes(apiRoute)

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
