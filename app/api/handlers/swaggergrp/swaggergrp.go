package swaggergrp

// Handlers manages the health endpoint used by k8s.
type Handlers struct {
}

// New constructs a Handlers type for route access.
func New() *Handlers {
	return &Handlers{}
}
