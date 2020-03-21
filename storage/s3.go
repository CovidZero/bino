package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gocloud.dev/blob"

	// suporte apenas para s3 por enquanto
	_ "gocloud.dev/blob/s3blob"
)

type (
	// S3Storage usa S3 para guardar os aruqivos temporários
	// Os buckets possuem o seguinte padrão /<nome da origem de dados>/<minuto da coleta>.data
	S3Storage struct {
		bucket *blob.Bucket
	}
)

// StoreCrawl guarda "data" em um bucket identificado pela origem e crawData
func (s *S3Storage) StoreCrawl(ctx context.Context, source string, crawlDate time.Time, format string, data []byte) (string, error) {
	crawlDate = crawlDate.Truncate(time.Minute)
	objName, err := ComputeObjectName(source, crawlDate, format)
	if err != nil {
		return "", err
	}
	w, err := s.bucket.NewWriter(ctx, objName, &blob.WriterOptions{
		// TODO: por hora, só JSON é usando, então isso arquivo resolve
		ContentType:        fmt.Sprintf("application/%s", format),
		ContentDisposition: fmt.Sprintf("rawData.%s", format),
		// TODO: checar se realmente é esse e se precisamos deixar isso configurável (aka passar como parâmetro na função)
		ContentEncoding: "utf-8",
	})
	if err != nil {
		return "", err
	}

	defer func() {
		// INFO: apenas para garantir que nenhum recurso fique aberto caso ocorra erro no envio dos dados
		if w != nil {
			if err := w.Close(); err != nil {
				log.Error().Err(err).Msg("Unexpected error closing bucket writer")
			}
		}
	}()

	if _, err = w.Write(data); err != nil {
		// INFO: evitando de empacotar o erro usando %w para não fazer a abstração de qual tipo de storage é usado
		return "", fmt.Errorf("unable to upload content to server, cause: %v", err)
	}
	err = w.Close()

	// INFO: evita de tentar chamar close duas vezes
	w = nil
	if err != nil {
		log.Error().
			Err(err).
			Str("module", "s3sstorage").
			Str("objectName", objName).
			Msg("Unable to close object. Data might be corrupted")
		return "", err
	}
	return objName, nil
}

// ComputeObjectName retorna o nome de um objeto na S3 que identifica uma coleta
func ComputeObjectName(source string, crawlDate time.Time, format string) (string, error) {
	if !validSource(source) {
		return "", errors.New("invalid source. please try again")
	}

	if !validCrawlDate(crawlDate) {
		return "", errors.New("invalid crawl data. precision is limited to the minute")
	}

	if !validFormat(format) {
		return "", errors.New("invalid format. please use json")
	}
	// INFO: formato deve ser <source>/<ano-mes-dia>/<hora-minuto>/rawData.<formato>
	return fmt.Sprintf("%s/%s/%s/rawData.%s", source, crawlDate.Format("2006-01-02"), crawlDate.Format("15-04"), format), nil
}

func validSource(source string) bool {
	return source == "ministerio_saude_brasil"
}

func validCrawlDate(crawlDate time.Time) bool {
	name, _ := crawlDate.Zone()
	return name == "UTC"
}

func validFormat(format string) bool {
	return format == "json"
}
