package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/temelpa/timetravel/service"
)

// GET /records/{id}
// GetRecord retrieves the record for the v2 endpoint
func (a *APIv2) GetRecordsV2(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	idNumber, err := strconv.ParseInt(id, 10, 32)

	if err != nil || idNumber <= 0 {
		err := writeError(w, "invalid id; id must be a positive number", http.StatusBadRequest)
		logError(err)
		return
	}

	record, err := service.GetRecordV2(
        ctx,
		a.db,
		int(idNumber),
	)
	if err != nil {
		err := writeError(w, fmt.Sprintf("record of id %v does not exist", idNumber), http.StatusBadRequest)
		logError(err)
		return
	}

	err = writeJSON(w, record, http.StatusOK)
	logError(err)
}

// GET /records/{id}
// GetRecord retrieves the record for the v1 endpoint
func (a *API) GetRecords(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	idNumber, err := strconv.ParseInt(id, 10, 32)

	if err != nil || idNumber <= 0 {
		err := writeError(w, "invalid id; id must be a positive number", http.StatusBadRequest)
		logError(err)
		return
	}

	record, err := a.records.GetRecord(
		ctx,
		int(idNumber),
	)
	if err != nil {
		err := writeError(w, fmt.Sprintf("record of id %v does not exist", idNumber), http.StatusBadRequest)
		logError(err)
		return
	}

	err = writeJSON(w, record, http.StatusOK)
	logError(err)
}
