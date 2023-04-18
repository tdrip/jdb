package jdb

import (
	"errors"
	"os"
)

type storage interface {
	ReadOnly() bool
	Exit() chan error
	Run(items map[string]IKeyedItem) (map[string]IKeyedItem, error)
}

func intiliase(filedb string, encdata EncodeKeyItems) error {
	if len(filedb) == 0 {
		return errors.New("file path for json database missing")
	}

	if _, err := os.ReadFile(filedb); err != nil {
		empty := make(map[string]IKeyedItem, 0)
		b, err := encdata(empty)
		if err != nil {
			return err
		}
		if err = os.WriteFile(filedb, b, 0644); err != nil {
			return err
		}
	}

	return nil
}

func processStorage(storage chan storage, db string, encdata EncodeKeyItems, decdata DecodeKeyItems) {
	for {
		s := <-storage

		content, err := os.ReadFile(db)
		if err != nil {
			s.Exit() <- err
			continue
		}

		converted, err := decdata(content)
		if err != nil {
			s.Exit() <- err
			continue
		}

		modified, err := s.Run(converted)
		if err != nil {
			s.Exit() <- err
			continue
		}
		if !s.ReadOnly() {
			if modified == nil {
				s.Exit() <- errors.New("modified data was nil")
				continue
			}

			b, err := encdata(modified)
			if err != nil {
				s.Exit() <- err
				continue
			}

			err = os.WriteFile(db, b, 0644)
			if err != nil {
				s.Exit() <- err
				continue
			}
		}
		s.Exit() <- err
	}
}
