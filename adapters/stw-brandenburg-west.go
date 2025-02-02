package adapters

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"mensa-queue/internal/payload"
	"mensa-queue/internal/repository"
)

type StwBrandenburgWestMensa struct {
	Provider     *StwBrandenburgWestAdapter
	StwID        int
	MensaHubID   *int32
	Uuid         string
	isRegistered bool
	Name         string
}

func (m *StwBrandenburgWestMensa) RegisterMensa(ctx context.Context, conn *pgx.Conn) (err error) {
	repo := repository.New(conn)
	id, err := repo.InsertMensa(ctx, repository.InsertMensaParams{
		Name:        m.Name,
		Description: nil,
		Slug:        m.Name,
		Uuid:        m.Uuid,
		ProviderID:  *m.Provider.MensaHubID,
	})
	if err != nil {
		return err
	}
	m.MensaHubID = &id
	m.isRegistered = true
	return nil
}

func (m *StwBrandenburgWestMensa) ParseMenu() ([]payload.LocalRecipe, error) {
	return nil, nil
}

func (m *StwBrandenburgWestMensa) IsRegistered() bool {
	return m.isRegistered
}

type StwBrandenburgWestAdapter struct {
	Name         string
	MensaHubID   *int32
	Uuid         string
	Mensas       []*StwBrandenburgWestMensa
	isRegistered bool
}

func NewStwBrandenburgWestAdapter(name string) *StwBrandenburgWestAdapter {
	adapter := &StwBrandenburgWestAdapter{
		Name: name,
		Uuid: "ad050c8a-cd7c-491c-a905-0aee84ae449a",
	}

	adapter.Mensas = []*StwBrandenburgWestMensa{
		{
			Provider: adapter,
			StwID:    9600,
			Name:     "Neues Palais",
			Uuid:     "7ee922a5-28d5-4726-a6a0-0905b79b3396",
		},
	}

	return adapter
}

func (a *StwBrandenburgWestAdapter) RegisterAdapter(ctx context.Context, conn *pgx.Conn) (err error) {
	repo := repository.New(conn)
	fmt.Println("Registering Adapter with UUID: ", a.Uuid)
	id, err := repo.InsertMensaProvider(ctx, repository.InsertMensaProviderParams{
		Name:        a.Name,
		Description: "Fill Me",
		Slug:        a.Name,
		Uuid:        a.Uuid,
	})
	if err != nil {
		return err
	}
	a.MensaHubID = &id
	a.isRegistered = true
	return nil
}

func (a *StwBrandenburgWestAdapter) GetAllMensas() []Mensa {
	mensas := make([]Mensa, len(a.Mensas))
	for i, mensa := range a.Mensas {
		mensas[i] = mensa // StwBrandenburgWestMensa implements Mensa, so this is valid
	}
	return mensas
}
