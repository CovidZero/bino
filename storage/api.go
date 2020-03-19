package storage

import "time"

type (
	// Temp define a API básica para guardar informações coletadas de modo que o processamento das mesmas
	// possa ser feito de forma independente, evitando sobrecarregar o servidor de origem (atualmente Ministério da Saúde)
	Temp interface {
		// StoreCrawl guarda uma coleta efetuada em uma origem no minuto informado
		StoreCrawl(source string, timeOfCrawl time.Time, resource []byte)
	}
)
