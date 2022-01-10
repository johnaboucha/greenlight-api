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
	Users  UserModel
}

// returns Models struct containing MovieModel and UserModel
func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
		Users:  UserModel{DB: db},
	}
}
