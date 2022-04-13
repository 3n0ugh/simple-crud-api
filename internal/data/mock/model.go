package mock

import "github.com/3n0ugh/simple-crud-api/internal/data"

func NewModel() data.Model {
	return data.Model{
		Book: BookModel{},
	}
}
