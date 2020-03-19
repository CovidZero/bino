package storage

import (
	"context"
	"time"
)

type (
	// Temp define a API básica para guardar informações coletadas de modo que o processamento das mesmas
	// possa ser feito de forma independente, evitando sobrecarregar o servidor de origem (atualmente Ministério da Saúde)
	Temp interface {
		// StoreCrawl guarda uma coleta efetuada em uma origem no minuto informado
		StoreCrawl(ctx context.Context, source string, timeOfCrawl time.Time, format string, resource []byte) (string, error)
	}
)
