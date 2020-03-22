package server

import (
	"net/http"
	"time"

	"github.com/CovidZero/bino/datasources"
	"github.com/CovidZero/bino/storage"
	"github.com/rs/zerolog/log"
)

type (
	// Crawl expões as operações de baixar e guardar arquivos em um storage temporário
	Crawl struct {
		source datasources.OnDemand
		temp   storage.Temp
	}
)

// FetchData inicia um novo processo de coleta (caso exista necessidade), caso contrário,
// retorna a última coleta feita.
func (c *Crawl) FetchData(w http.ResponseWriter, req *http.Request) {
	now := time.Now().UTC().Truncate(time.Minute)
	buf, err := c.source.Collect(now)
	if err != nil {
		log.Error().Err(err).Str("module", "crawl").Str("source", c.source.Name()).Msg("Error reading data from source")
		http.Error(w, "Unexpected error. Try again later", http.StatusBadGateway)
		return
	}

	name, err := c.temp.StoreCrawl(req.Context(), "ministerio_saude_brasil", now, c.source.Format().Ext(), buf)
	if err != nil {
		log.Error().Err(err).Str("module", "crawl").Msg("Unable to send data to temporary db")
		http.Error(w, "Unexpected error. Try again later", http.StatusBadGateway)
		return
	}
	respondWithJSON(w, req, struct {
		Path string `json:"path"`
	}{Path: name})
}
