package jdb

type IKeyedItem interface {
	// get a unique for the keyed item
	GetID() string

	// generate a new ID for the item
	NewID() string
}
