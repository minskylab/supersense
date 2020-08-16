package server

import (
	"github.com/minskylab/supersense"
	"github.com/minskylab/supersense/persistence"
)

func New(mux *supersense.Mux, db *persistence.Persistence) (*Server, error) {
	return &Server{
		mux: mux,
		db:  db,
	}, nil
}