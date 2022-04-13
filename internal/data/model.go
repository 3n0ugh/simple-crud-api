package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Model struct {
	Book interface {
		Insert(book *Book) error
		GetAll() ([]*Book, error)
		Update(book *Book) error
		Delete(id int64) error
	}
}

func NewModel(db *sql.DB) Model {
	return Model{
		Book: &BookModel{DB: db},
	}
}
