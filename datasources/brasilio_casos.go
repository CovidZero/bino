package datasources

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
)

type (
	// BrasilIOCovid19Dataset é um datasource OnDemand que retorna os dados do site https://brasil.io/dataset/covid19
	BrasilIOCovid19Dataset struct {
		// Endpoint contém a url para ser acionada, por padrão utiliza
		// https://brasil.io/dataset/covid19/caso?format=csv
		endpoint *url.URL
	}
)

const (
	// BrasilIOCovid19DefaultURL é a URL padrão a ser utilizada caso a variável DATASOURCE_IVIS_URL não exista
	BrasilIOCovid19DefaultURL = "https://brasil.io/dataset/covid19/caso?format=csv"

	// BrasilIOCovid19DatasetName é o nome pelo qual o dataset será referenciado na S3
	BrasilIOCovid19DatasetName = "brasil_io_covid19"
)

var (
	defaultBrasilIOCovid19 = &IVISDataset{}
)

// Name identificar o dataset
func (i *BrasilIOCovid19Dataset) Name() string { return BrasilIOCovid19DatasetName }

// Format retorna o formato no qual os dados são coletados
func (i *BrasilIOCovid19Dataset) Format() Format {
	return CSV
}

// Encoding retorna a codificação dos dados
func (i *BrasilIOCovid19Dataset) Encoding() string {
	// TODO: checar se essa codificação está correta
	return "utf-8"
}

// Collect aciona o serviço e roda uma coleta
func (i *BrasilIOCovid19Dataset) Collect(now time.Time, args url.Values) ([]byte, error) {
	// TODO: gerar os logs aqui é incorreto pois isso é responsabilidade de quem chamou, mas por hora serve
	// TODO: atlerar o http client e incluir um timeout para proteger o nosso processo caso o servidor remoto esteja muito lento
	var err error
	var actualEndpoint url.URL
	if i.endpoint != nil {
		actualEndpoint = *i.endpoint
	} else {
		u, err := url.Parse(BrasilIOCovid19DefaultURL)
		if err != nil {
			return nil, err
		}
		actualEndpoint = *u
	}

	params := actualEndpoint.Query()
	for k, values := range args {
		for _, v := range values {
			params.Add(k, v)
		}
	}
	actualEndpoint.RawQuery = params.Encode()

	log.Info().Str("module", "datasource").Str("crawler", "BrasilIOCovid19").Str("url", actualEndpoint.String()).Send()
	response, err := http.Get(actualEndpoint.String())
	if err != nil {
		// TODO: 500 nesse caso não é o ideal mas por enquanto resolve
		log.Error().Err(err).Str("module", "datasource").Str("crawler", "BrasilIOCovid19").Msg("Error reaching source")
		return nil, err
	}
	if response.StatusCode >= 400 {
		log.Error().Err(errors.New(response.Status)).
			Int("status", response.StatusCode).
			Str("module", "datasource").
			Str("crawler", "BrasilIOCovid19").
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
			Str("crawler", "BrasilIOCovid19").
			Msg("Unable to read response body")
		return nil, err
	}
	return buf, nil
}
