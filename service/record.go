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
    // RecordService method for v1 endpoint

	// GetRecord will retrieve an record.
	GetRecord(ctx context.Context, id int) (entity.Record, error)

	// CreateRecord will insert a new record.
	//
	// If it a record with that id already exists it will fail.
	CreateRecord(ctx context.Context, record entity.Record) error

	// UpdateRecord will change the internal `Map` values of the record if they exist.
	// if the update[key] is null it will delete that key from the record's Map.
	//
	// UpdateRecord will error if id <= 0 or the record does not exist with that id.
	UpdateRecord(ctx context.Context, id int, updates map[string]*string) (entity.Record, error)
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

func (s *InMemoryRecordService) GetRecord(ctx context.Context, id int) (entity.Record, error) {
	record := s.data[id]
	if record.ID == 0 {
		return entity.Record{}, ErrRecordDoesNotExist
	}

	record = record.Copy() // copy is necessary so modifations to the record don't change the stored record
	return record, nil
}

func (s *InMemoryRecordService) CreateRecord(ctx context.Context, record entity.Record) error {
	id := record.ID
	if id <= 0 {
		return ErrRecordIDInvalid
	}

	existingRecord := s.data[id]
	if existingRecord.ID != 0 {
		return ErrRecordAlreadyExists
	}

	s.data[id] = record
	return nil
}

func (s *InMemoryRecordService) UpdateRecord(ctx context.Context, id int, updates map[string]*string) (entity.Record, error) {
	entry := s.data[id]
	if entry.ID == 0 {
		return entity.Record{}, ErrRecordDoesNotExist
	}

	for key, value := range updates {
		if value == nil { // deletion update
			delete(entry.Data, key)
		} else {
			entry.Data[key] = *value
		}
	}

	return entry.Copy(), nil
}

// GetRecordV2 serves the v2 endpoint
// will retrieve all records with the passed in id value (rid)
// Searches the records table for all records for a given rid value
// Returns a map of the most up to date assignment of (key,value) pairs
func GetRecordV2(ctx context.Context, db *sql.DB, id int) (entity.Record, error) {
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

// AddRecordRowV2 serves the v2 endpoint
// Inserts a new record
// rid param is the record id which is the record that is being added to (X in /api/v1/records/X)
// eid is a counter for each id in the records sql database (not returned to the user)
// key, value are added to the records table
func AddRecordRowV2(ctx context.Context, db *sql.DB, rid int, eid int, key string, value *string) error {
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

