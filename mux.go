package supersense

// Mux is a necesary struct to join different sources
type Mux struct {
	sources []Source
}

// NewMux returns a new mux
func NewMux(sources ...Source) (*Mux, error) {
	return &Mux{sources: sources}, nil
}
