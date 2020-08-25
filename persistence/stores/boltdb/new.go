package boltdb

import (
	"github.com/asdine/storm/v3"
	"github.com/pkg/errors"
)

func NewStore(path string, maxBuffer int64) (*Store, error) {
	db, err := storm.Open(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if maxBuffer <= 0 {
		maxBuffer = 100e3
	}
	s := &Store{db: db, mainID: mainIDValue, maxBuffer: maxBuffer}

	if err = s.initialCheck(); err != nil {
		return nil, errors.WithStack(err)
	}

	return s, nil
}
