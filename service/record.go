package service

import (
	"context"
	"errors"
    "database/sql"

	"github.com/temelpa/timetravel/entity"
)

var ErrRecordDoesNotExist = errors.New("record with that id does not exist")
var ErrRecordIDInvalid = errors.New("record id must >= 0")
var ErrRecordAlreadyExists = errors.New("record already exists")

// Implements method to get, create, and update record data.
type RecordService interface {

	// GetRecord will retrieve all records with the passed in id value (rid)
    // Searches the records table for all records for a given rid value
    // Returns a map of the most up to date assignment of (key,value) pairs
	GetRecord(ctx context.Context, id int) (entity.Record, error)

	// AddRecordRow will insert a new record.
    //
    // rid param is the record id which is the record that is being added to
    // eid is a counter for each id in the records table
    // key, value are added to the records table
	AddRecordRow(ctx context.Context, db *sql.DB, rid int, eid int, key string, value *string) error
}

// InMemoryRecordService is an in-memory implementation of RecordService.
type InMemoryRecordService struct {
	data map[int]entity.Record
}

func NewInMemoryRecordService() InMemoryRecordService {
	return InMemoryRecordService{
		data: map[int]entity.Record{},
	}
}
func GetRecord(ctx context.Context, db *sql.DB, id int) (entity.Record, error) {
    rows, err := db.Query("SELECT id, rid, key, value FROM records WHERE rid=? ORDER BY id DESC", id)
    if err != nil {
        return entity.Record{}, err
    }

    record := entity.Record{}
    record.ID = id
    record.Data = map[string]string{}
    for rows.Next() {
        r := entity.RecordRow{}
        err := rows.Scan(&r.ID, &r.RID, &r.Key, &r.Value)

        if err != nil {
            return entity.Record{}, err
        }
        _, ok := record.Data[r.Key]
        // Only store the last value for a (key, value) pair since the Query orders the
        // results by last update entry first
        if !ok {
            if r.Value.Valid {
                record.Data[r.Key] = r.Value.String
            } else {
                record.Data[r.Key] = ""
            }
        }
    }

    if len(record.Data) == 0 {
        return entity.Record{}, ErrRecordDoesNotExist
    }

    // Delete all the entries with null values in the response
    for key, value := range record.Data {
        if value == "" {
            delete(record.Data, key)
        }
    }

	return record, nil
}

func AddRecordRow(ctx context.Context, db *sql.DB, rid int, eid int, key string, value *string) error {
    // rid is the id for the record (X in /api/v1/records/X)
    // eid is the id for the individual entry within the db (this does not get returned to the user)
    stmt, err := db.Prepare("INSERT INTO records VALUES(?,?,?,?);")
    if err != nil {
        return err
    }
    res, err := stmt.Exec(eid, rid, key, value)

    if err != nil {
        return err
    }
    if _, err := res.LastInsertId(); err != nil {
        return err
    }

    stmt.Close()
    return nil
}

