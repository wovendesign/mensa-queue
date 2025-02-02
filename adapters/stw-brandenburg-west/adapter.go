package stw_brandenburg_west

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"mensa-queue/adapters"
	"mensa-queue/internal/repository"
)

type StwBrandenburgWestAdapter struct {
	Name         string
	MensaHubID   *int32
	Uuid         string
	Mensas       []*StwBrandenburgWestMensa
	isRegistered bool
}

func NewAdapter(name string) *StwBrandenburgWestAdapter {
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
		{
			Provider: adapter,
			StwID:    9601,
			Name:     "Griebnitzsee",
			Uuid:     "d88eff2d-2f3f-4e03-a77d-a035ffeda4d2",
		},
		{
			Provider: adapter,
			StwID:    9602,
			Name:     "Golm",
			Uuid:     "097826b0-3943-4727-8307-649179bc1d76",
		},
		{
			Provider: adapter,
			StwID:    9603,
			Name:     "Filmuniversität",
			Uuid:     "5bfef986-0531-4fc5-a5e8-bae7dfbf8c20",
		},
		{
			Provider: adapter,
			StwID:    9604,
			Name:     "FHP",
			Uuid:     "799301f4-efa3-41cf-914b-c2b1ed58bb1f",
		},
		{
			Provider: adapter,
			StwID:    9605,
			Name:     "Wildau",
			Uuid:     "52f71c86-5e10-42bf-8b23-64d45fe62fa5",
		},
		{
			Provider: adapter,
			StwID:    9606,
			Name:     "Brandenburg",
			Uuid:     "9c09d17d-57dc-43f1-892c-ac11e03f492f",
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

func (a *StwBrandenburgWestAdapter) GetAllMensas() []adapters.Mensa {
	mensas := make([]adapters.Mensa, len(a.Mensas))
	for i, mensa := range a.Mensas {
		mensas[i] = mensa // StwBrandenburgWestMensa implements Mensa, so this is valid
	}
	return mensas
}
