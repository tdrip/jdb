package jdb

import (
	"fmt"
)

// read one item from the database
type readStorageItem struct {
	id   string
	item chan IKeyedItem
	exit chan error
}

func newReadStorageItem(id string) *readStorageItem {
	return &readStorageItem{
		id:   id,
		item: make(chan IKeyedItem, 1),
		exit: make(chan error, 1),
	}
}

func (rsi readStorageItem) ReadOnly() bool {
	return true
}

func (rsi readStorageItem) Exit() chan error {
	return rsi.exit
}

func (rsi readStorageItem) Run(items map[string]IKeyedItem) (map[string]IKeyedItem, error) {
	_, exists := items[rsi.id]

	if !exists {
		err := fmt.Errorf("id %s does not exist", rsi.id)
		return nil, err
	}

	rsi.item <- items[rsi.id]
	return items, nil
}
