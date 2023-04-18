package jdb

import (
	"errors"
)

// save one item to the database
type saveStorageItem struct {
	item  IKeyedItem
	saved chan IKeyedItem
	exit  chan error
}

func newSaveStorageItem(item IKeyedItem) *saveStorageItem {
	return &saveStorageItem{
		item:  item,
		saved: make(chan IKeyedItem, 1),
		exit:  make(chan error, 1),
	}
}

func (ssi saveStorageItem) ReadOnly() bool {
	return false
}

func (ssi saveStorageItem) Exit() chan error {
	return ssi.exit
}

func (ssi saveStorageItem) Run(items map[string]IKeyedItem) (map[string]IKeyedItem, error) {
	item := ssi.item
	if item == nil {
		err := errors.New("item to be saved was nil")
		return nil, err
	}

	id := item.GetID()
	if len(id) == 0 {
		err := errors.New("item id was empty")
		return nil, err
	}

	items[id] = item
	ssi.saved <- item
	return items, nil
}
