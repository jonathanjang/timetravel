package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/temelpa/timetravel/entity"
	"github.com/temelpa/timetravel/service"
)

// POST /records/{id}
// if the record exists, the record is updated.
// if the record doesn't exist, the record is created.
func (a *API) PostRecords(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	idNumber, err := strconv.ParseInt(id, 10, 32)

	if err != nil || idNumber <= 0 {
		err := writeError(w, "invalid id; id must be a positive number", http.StatusBadRequest)
		logError(err)
		return
	}

	var body map[string]*string
	err = json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		err := writeError(w, "invalid input; could not parse json", http.StatusBadRequest)
		logError(err)
		return
	}

	// first retrieve the record
    getRecord, _ := service.GetRecord(
		ctx,
        a.db,
		int(idNumber),
	)

    responseRecord := entity.Record{}
    responseRecord.ID = int(idNumber)
    responseRecord.Data = map[string]string{}

    // Start by pre-adding all entries from the GET Query to response
    for key, value := range getRecord.Data {
        responseRecord.Data[ key ] = value
    }

    for key, value := range body {
        a.IncrementEid(); // increment counter for Id in db
        err = service.AddRecordRow(
            ctx,
            a.db,
            int(idNumber),
            a.eid,
            key,
            value,
        )
        // Don't add null values to the response
        if value != nil {
            responseRecord.Data[key] = *value
        }
        // if the value changed to null, remove it from the response
        if _,ok := responseRecord.Data[key]; ok && value == nil {
            delete(responseRecord.Data, key)
        } else if value != nil {
            responseRecord.Data[key] = *value
        }
        if err != nil {
            errInWriting := writeError(w, ErrInternal.Error(), http.StatusInternalServerError)
	        logError(err)
	        logError(errInWriting)
		    return
	    }
	}
    err = writeJSON(w, responseRecord, http.StatusOK)
	logError(err)
}
