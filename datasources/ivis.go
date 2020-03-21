package datasources

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type (
	// IVISDataset é um datasource OnDemand que retorna os dados do ministério da saúde
	IVISDataset struct {
		// Endpoint contém a url para ser acionada, por padrão utiliza
		// http://plataforma.saude.gov.br/novocoronavirus/resources/scripts/database.js
		endpoint string
	}
)

const (
	IVISDatasetDefaultURL = "http://plataforma.saude.gov.br/novocoronavirus/resources/scripts/database.js"
)

var (
	defaultIVIS = &IVISDataset{}
)

// Name identificar o dataset
func (i *IVISDataset) Name() string { return "ministerio_saude_brasil" }

// Format retorna o formato no qual os dados são coletados
func (i *IVISDataset) Format() Format {
	return JSON
}

// Encoding retorna a codificação dos dados
func (i *IVISDataset) Encoding() string {
	// TODO: checar se essa codificação está correta
	return "iso-8859-1"
}

// Collect aciona o serviço e roda uma coleta
func (i *IVISDataset) Collect(_ time.Time) ([]byte, error) {
	// TODO: gerar os logs aqui é incorreto pois isso é responsabilidade de quem chamou, mas por hora serve
	// TODO: checar se podemos trocar por https para garantir integridade dos dados
	// TODO: atlerar o http client e incluir um timeout para proteger o nosso processo caso o servidor remoto esteja muito lento
	url := i.endpoint
	if url == "" {
		url = IVISDatasetDefaultURL
	}
	response, err := http.Get(url)
	if err != nil {
		// TODO: 500 nesse caso não é o ideal mas por enquanto resolve
		log.Error().Err(err).Str("module", "datasource").Str("crawler", "IVISDataset").Msg("Error reaching source")
		return nil, err
	}
	if response.StatusCode >= 400 {
		log.Error().Err(errors.New(response.Status)).
			Int("status", response.StatusCode).
			Str("module", "datasource").
			Str("crawler", "IVISDataset").
			Msg("Unexpected status from source")
		return nil, err
	}
	response.Close = true
	defer response.Body.Close()

	// TODO vetor de ataque em potencial, pois não existe limite no tamanho do buffer que vai ficar em memória (melhorar isso)
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error().Err(errors.New(response.Status)).
			Int("status", response.StatusCode).
			Str("module", "datasource").
			Str("crawler", "IVISDataset").
			Msg("Unable to read response body")
		return nil, err
	}

	if bytes.HasPrefix(buf, []byte("var database=")) {
		buf = buf[len("var database="):]
	}
	return buf, nil
}
