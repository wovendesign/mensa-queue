package stw_brandenburg_west

import (
	"context"
	"github.com/jackc/pgx/v5"
	"mensa-queue/internal/repository"
)

type StwBrandenburgWestMensa struct {
	Provider     *StwBrandenburgWestAdapter
	StwID        int32
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

func (m *StwBrandenburgWestMensa) IsRegistered() bool {
	return m.isRegistered
}

func (m *StwBrandenburgWestMensa) AiGenerationEnabled() bool {
	return false
}
