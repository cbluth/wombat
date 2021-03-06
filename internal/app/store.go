package app

import (
	"path/filepath"

	badger "github.com/dgraph-io/badger/v2"
)

type dblogger = badger.Logger

type store struct {
	db *badger.DB
}

func newStore(path string, l dblogger) (*store, error) {
	dbPath := filepath.Join(path, "db")
	opts := badger.DefaultOptions(dbPath)
	opts = opts.WithLogger(l)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &store{db}, nil
}

func (s *store) get(key []byte) (val []byte, rtnErr error) {
	rtnErr = s.db.View(func(txn *badger.Txn) error {
		data, err := txn.Get(key)
		if err != nil {
			return err
		}
		val, err = data.ValueCopy(val)
		return err
	})
	return
}

func (s *store) set(key, val []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, val)
	})
}

func (s *store) close() {
	if s == nil {
		return
	}
	s.db.Close()
}
