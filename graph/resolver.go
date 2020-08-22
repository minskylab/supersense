package graph

import (
	"sync"

	"github.com/minskylab/supersense"
	"github.com/minskylab/supersense/sources"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver ...
type Resolver struct {
	mu  *sync.Mutex
	mux *supersense.Mux
	spokesman *sources.Spokesman
}

// NewResolver returns a new resolver instance
func NewResolver(mux *supersense.Mux, spokesman *sources.Spokesman) *Resolver {
	return &Resolver{mux: mux, mu: &sync.Mutex{}, spokesman: spokesman}
}
