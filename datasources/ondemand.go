package datasources

import "errors"

var (
	errMissingDatasource = errors.New("datasource name is not valid")
)

// GetOnDemand retorna o datasorce identificado pelo nome ou um erro caso o nome seja inválido
func GetOnDemand(name string) (OnDemand, error) {
	// TODO: pensar em uma forma de registrar os datasources, por hora isso arquivo resolve
	switch name {
	case "ministerio_saude_brasil":
		return &IVISDataset{}, nil
	}
	return nil, errMissingDatasource
}
