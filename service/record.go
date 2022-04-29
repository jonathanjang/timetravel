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
	GetRecord(ctx context.Context, id int) (entity.Record, error)

	// AddRecordRow will insert a new record.
    //
    // rid param is the record id which is the record that is being added to
    // eid is a counter for each id in the records table
    // key, value are added to the records table
	AddRecordRow(ctx context.Context, db *sql.DB, rid int, eid int, key string, value *string) error

    // UpdateOrDeleteRecord takes in the recordID, key, and value and checks if
    // an update or delete is required. If the value field is null, the entry
    // in the database is deleted. Otherwise, an update is performed on the data
	UpdateOrDeleteRecord(ctx context.Context, db *sql.DB, rid int, key string, prevValue string,
                 newValue *string) error
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
// TODO: match error checking with non sql skeleton
func GetRecord(ctx context.Context, db *sql.DB, id int) (entity.Record, error) {
    rows, err := db.Query("SELECT rid, key, value FROM records WHERE rid=?", id)
    if err != nil {
        return entity.Record{}, err
    }

    record := entity.Record{}
    record.ID = id
    record.Data = map[string]string{}
    for rows.Next() {
        r := entity.RecordRow{}
        err := rows.Scan(&r.ID, &r.Key, &r.Value)
        if err != nil {
            return entity.Record{}, err
        }
        record.Data[r.Key] = r.Value
    }

    if len(record.Data) == 0 {
        return entity.Record{}, ErrRecordDoesNotExist
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

func UpdateOrDeleteRecord(ctx context.Context, db *sql.DB, rid int, key string, prevValue string, newValue *string) error {
    if newValue != nil {
        stmt, err := db.Prepare("UPDATE records set key=?, value=? WHERE rid=? AND key=? AND value=?")
        if err != nil {
            return err
        }
        _, err = stmt.Exec(key, newValue, rid, key, prevValue)
        if err != nil {
            return err
        }

        stmt.Close()
    } else {
        stmt, err := db.Prepare("DELETE from records WHERE rid=? AND key=? AND value=?")
        if err != nil {
            return err
        }
        _, err = stmt.Exec(rid, key, prevValue)
        if err != nil {
            return err
        }

        stmt.Close()
    }
    return nil
}
