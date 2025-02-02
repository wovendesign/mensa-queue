package adapters

import (
	"context"
	"github.com/jackc/pgx/v5"
	"mensa-queue/models"
)

// MenuParser is the common interface for all canteen adapters.
type Adapter interface {
	RegisterAdapter(ctx context.Context, conn *pgx.Conn) (err error)
	GetAllMensas() []Mensa
}

type Mensa interface {
	ParseMenu() ([]*models.Recipe, error)
	RegisterMensa(ctx context.Context, conn *pgx.Conn) (err error)
	IsRegistered() bool
	AiGenerationEnabled() bool
}
