package adapters

import "mensa-queue/models"

type StwBrandenburgWestMensa struct {
	StwID      uint
	MensaHubID int8
	Name       string
}

func (m *StwBrandenburgWestMensa) RegisterMensa() {
	// make request to db, check if the mensa is registered, if not, register it

	// then, set the internal MensaHubID of the mensa
	m.MensaHubID = 1
}

func (m *StwBrandenburgWestMensa) ParseMenu() ([]models.MenuItem, error) {
	return nil, nil
}

type StwBrandenburgWestAdapter struct {
	Name   string
	Mensas []*StwBrandenburgWestMensa
}

func NewStwBrandenburgWestAdapter(name string) *StwBrandenburgWestAdapter {
	return &StwBrandenburgWestAdapter{
		Name: name,
		Mensas: []*StwBrandenburgWestMensa{
			{
				StwID: 9600,
				Name:  "Neues Palais",
			},
		},
	}
}

func (a *StwBrandenburgWestAdapter) RegisterAdapter() {

}

func (a *StwBrandenburgWestAdapter) GetAllMensas() []Mensa {
	mensas := make([]Mensa, len(a.Mensas))
	for i, mensa := range a.Mensas {
		mensas[i] = mensa // StwBrandenburgWestMensa implements Mensa, so this is valid
	}
	return mensas
}
