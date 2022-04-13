package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/3n0ugh/simple-crud-api/internal/validator"
	"time"
)

type Book struct {
	ID     int64  `json:"id,omitempty"`
	Name   string `json:"name"`
	Author string `json:"author"`
}

type BookModel struct {
	DB *sql.DB
}

func ValidateBook(v *validator.Validator, book *Book) {
	v.Check(book.Name != "", "name", "must be provided")
	v.Check(len(book.Name) <= 40, "name", "must not be more than 40 bytes long")

	v.Check(book.Author != "", "author", "must be provided")
	v.Check(len(book.Author) <= 40, "author", "must not be more than 40 bytes long")
}

func (b *BookModel) Insert(book *Book) error {
	query := `INSERT INTO books(name, author) VALUES($1,$2) RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{book.Name, book.Author}

	return b.DB.QueryRowContext(ctx, query, args...).
		Scan(&book.ID)
}

func (b *BookModel) GetAll() ([]*Book, error) {
	query := `SELECT id, name, author FROM books`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := b.DB.QueryContext(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	defer rows.Close()

	var books = make([]*Book, 0)

	for rows.Next() {
		var book Book

		err = rows.Scan(&book.ID, &book.Name, &book.Author)
		if err != nil {
			return nil, err
		}

		books = append(books, &book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (b *BookModel) Update(book *Book) error {
	query := `UPDATE books SET name=$1, author=$2 WHERE id=$3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{book.Name, book.Author, book.ID}

	res, err := b.DB.ExecContext(ctx, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}

	if r, _ := res.RowsAffected(); r == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (b *BookModel) Delete(id int64) error {
	query := `DELETE FROM books WHERE id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := b.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	effectedRow, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if effectedRow == 0 {
		return ErrRecordNotFound
	}

	return nil
}
