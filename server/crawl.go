package server

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/CovidZero/bino/storage"
	"github.com/rs/zerolog/log"
)

type (
	// Crawl expões as operações de baixar e guardar arquivos em um storage temporário
	Crawl struct {
		temp storage.Temp
	}
)

// FetchData inicia um novo processo de coleta (caso exista necessidade), caso contrário,
// retorna a última coleta feita
func (c *Crawl) FetchData(w http.ResponseWriter, req *http.Request) {
	// TODO: checar se podemos trocar por https para garantir integridade dos dados
	response, err := http.Get("http://plataforma.saude.gov.br/novocoronavirus/resources/scripts/database.js")
	if err != nil {
		// TODO: 500 nesse caso não é o ideal mas por enquanto resolve
		log.Error().Err(err).Str("module", "crawl").Msg("Error reaching source")
		http.Error(w, "Unable to download content", http.StatusInternalServerError)
		return
	}
	if response.StatusCode >= 400 {
		log.Error().Err(errors.New(response.Status)).Int("status", response.StatusCode).Str("module", "crawl").Msg("Unexpected status from source")
		http.Error(w, "Unable to download content. Server returned unexpected status", http.StatusInternalServerError)
		return
	}
	response.Close = true
	defer response.Body.Close()
	// TODO vetor de ataque em potencial, pois não existe limite no tamanho do buffer que vai ficar em memória (melhorar isso)
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error().Err(err).Str("module", "crawl").Msg("Error reading response body")
		http.Error(w, "Unable to download content. Server returned unexpected status", http.StatusInternalServerError)
		return
	}

	if bytes.HasPrefix(buf, []byte("var database=")) {
		buf = buf[len("var database="):]
	}

	name, err := c.temp.StoreCrawl(req.Context(), "ministerio_saude_brasil", time.Now().UTC(), "json", buf)
	if err != nil {
		log.Error().Err(err).Str("module", "crawl").Msg("Unable to send data to temporary db")
		return
	}
	respondWithJSON(w, req, struct {
		Path string `json:"path"`
	}{Path: name})
}
