package datasources

import (
	"errors"
	"net/url"
	"os"
)

var (
	errMissingDatasource = errors.New("datasource name is not valid")
)

// AllOnDemand retorna a lista de todos os coletores OnDemand disponíveis
func AllOnDemand() []string {
	return []string{
		IVISDatasetName,
		BrasilIOCovid19DatasetName,
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
	case BrasilIOCovid19DatasetName:
		endpoint := os.Getenv("DATASOURCE_BRASIL_IO_COVID19_URL")
		if endpoint == "" {
			return &BrasilIOCovid19Dataset{}, nil
		}
		u, err := url.Parse(endpoint)
		if err != nil {
			return nil, err
		}
		return &BrasilIOCovid19Dataset{
			endpoint: u,
		}, nil
	}
	return nil, errMissingDatasource
}
