package mock

import "github.com/3n0ugh/simple-crud-api/internal/data"

var Book = &data.Book{
	ID:     1,
	Name:   "Foo",
	Author: "Bar",
}

var Books = []*data.Book{
	{
		ID:     1,
		Name:   "Foo",
		Author: "Bar",
	},
	{
		ID:     2,
		Name:   "Lorem",
		Author: "Ipsum",
	},
}

type BookModel struct{}

func (b BookModel) Insert(book *data.Book) error {
	book.ID = Book.ID
	return nil
}
func (b BookModel) GetAll() ([]*data.Book, error) {
	return Books, nil
}
func (b BookModel) Update(book *data.Book) error {
	if book.ID != Book.ID {
		return data.ErrEditConflict
	}
	return nil
}
func (b BookModel) Delete(id int64) error {
	if id == Book.ID {
		return nil
	}
	return data.ErrRecordNotFound
}
