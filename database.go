package jdb

type EncodeKeyItems func(map[string]IKeyedItem) ([]byte, error)
type DecodeKeyItems func([]byte) (map[string]IKeyedItem, error)

type Database struct {
	Storage chan storage
}

func BuildDatabase(filedb string, encdata EncodeKeyItems, decdata DecodeKeyItems) (*Database, error) {

	err := intiliase(filedb, encdata)
	if err != nil {
		return nil, err
	}

	// create storage channel to communicate over
	store := make(chan storage)

	// start watching storage channel for work
	go processStorage(store, filedb, encdata, decdata)

	db := &Database{Storage: store}

	return db, nil
}

func (d *Database) SaveItem(todo IKeyedItem) (IKeyedItem, error) {
	job := newSaveStorageItem(todo)
	d.Storage <- job
	if err := <-job.Exit(); err != nil {
		return nil, err
	}
	return <-job.saved, nil
}

func (d *Database) CountItems() (int, error) {
	job := newCountStorageItems()
	d.Storage <- job

	if err := <-job.Exit(); err != nil {
		return -1, err
	}

	total := <-job.total
	return total, nil
}

func (d *Database) GetItems() ([]IKeyedItem, error) {
	arr := make([]IKeyedItem, 0)
	job := newReadStorageItems()
	d.Storage <- job

	if err := <-job.Exit(); err != nil {
		return arr, err
	}

	items := <-job.items
	for _, value := range items {
		arr = append(arr, value)
	}
	return arr, nil
}

func (d *Database) GetItem(id string) (IKeyedItem, error) {
	job := newReadStorageItem(id)
	d.Storage <- job
	if err := <-job.Exit(); err != nil {
		return nil, err
	}
	item := <-job.item
	return item, nil
}

func (d *Database) DeleteItem(id string) error {
	job := newDeleteStorageItem(id)
	d.Storage <- job

	if err := <-job.Exit(); err != nil {
		return err
	}
	return nil
}
