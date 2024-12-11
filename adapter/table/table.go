package table

import "errors"

type Table interface {
	Insert(item string) error
}

type table struct {
	items map[string]*struct{}
}

func New() *table {
	return &table{items: make(map[string]*struct{})}
}

func (t *table) Insert(item string) error {
	if t.items[item] != nil {
		return errors.New("item already exists in table")
	}

	t.items[item] = &struct{}{}
	return nil
}
