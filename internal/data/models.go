package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found") // not found errors
	ErrEditConflict   = errors.New("edit conflict")    // used when race conditions occur
)

// Models struct wraps our models
type Models struct {
	Movies MovieModel
}

// returns Models struct containing MovieModel
func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
	}
}
