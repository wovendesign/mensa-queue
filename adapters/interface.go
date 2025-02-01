package adapters

import "mensa-queue/models"

// MenuParser is the common interface for all canteen adapters.
type Adapter interface {
	RegisterAdapter()
	GetAllMensas() []Mensa
}

type Mensa interface {
	ParseMenu() ([]models.MenuItem, error)
	RegisterMensa()
}
