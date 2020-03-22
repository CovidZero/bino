package datasources

import "time"

type (
	// OnDemand indica data sources que retornam os dados em batches, normalmente não são atualizados com grande
	// frequencia e portanto não são coletados com frequência
	OnDemand interface {
		// Name indica o nome do datasource, deve ser um identificador válido (aka letras/números/sem espaços)
		Name() string

		// Format retorna o formato no qual os dados são gravados
		Format() Format

		// Encoding retorna a codificação usada pelos dados
		Encoding() string

		// Collect executa uma coleta de dados e retorna quando a coleta estiver concluída
		Collect(time.Time) ([]byte, error)
	}
)
