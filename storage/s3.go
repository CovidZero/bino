package storage

import (
	"time"
)

type (
	// S3Storage usa S3 para guardar os aruqivos temporários
	// Os buckets possuem o seguinte padrão /<nome da origem de dados>/<minuto da coleta>.data
	S3Storage struct {
		// TODO: actually implement this
	}
)

// StoreCrawl guarda "data" em um bucket identificado pela origem e crawData
func (s *S3Storage) StoreCrawl(source string, crawlData time.Time, data []byte) {

}
