package datasources

import (
	"errors"
	"os"
)

var (
	errMissingDatasource = errors.New("datasource name is not valid")
)

// AllOnDemand retorna a lista de todos os coletores OnDemand disponíveis
func AllOnDemand() []string {
	return []string{
		IVISDatasetName,
	}
}

// GetOnDemand retorna o datasorce identificado pelo nome ou um erro caso o nome seja inválido
func GetOnDemand(name string) (OnDemand, error) {
	// TODO: pensar em uma forma de registrar os datasources, por hora isso arquivo resolve
	switch name {
	case IVISDatasetName:
		return &IVISDataset{
			endpoint: os.Getenv("DATASOURCE_IVIS_URL"),
		}, nil
	}
	return nil, errMissingDatasource
}
