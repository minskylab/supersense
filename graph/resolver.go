package graph

import "github.com/minskylab/supersense"

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver ...
type Resolver struct {
	mux *supersense.Mux
}

// NewResolver returns a new resolver instance
func NewResolver(mux *supersense.Mux) *Resolver {
	return &Resolver{mux: mux}
}
