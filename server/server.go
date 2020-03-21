package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/CovidZero/bino/datasources"
	"github.com/CovidZero/bino/storage"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func allCollectors(storage storage.Temp) (http.Handler, error) {
	r := mux.NewRouter()
	err := setupDatasourcesRoutes(r, storage)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// NewAPI retorna um servidor HTTP pré-configurado que guarda as informações temporárias no storage informado
func NewAPI(bindAddr string, storage storage.Temp) (*http.Server, error) {
	collectors, err := allCollectors(storage)
	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr: bindAddr,
		// Proteger o server de clients muito lentos, evitando que existam muitas conexões TCP abertas ao mesmo tempo
		// IMPORTANTE! não protege de clientes que intencionamente enviam dados a uma velocidade muito lenta (Slow Loris Attack)
		// https://www.youtube.com/watch?v=XiFkyR35v2Y para mais informações
		ReadTimeout: time.Second * 10,

		// Mesma lógica acima e mesma limitação
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 * 1024 * 1024 * 1024,

		Handler: collectors,
	}
	return server, nil
}

func setupDatasourcesRoutes(r *mux.Router, storage storage.Temp) error {
	if err := registerCrawlEndpoint(r, storage, "ministerio_saude_brasil"); err != nil {
		return err
	}
	return nil
}

func registerCrawlEndpoint(r *mux.Router, storage storage.Temp, name string) error {
	source, err := datasources.GetOnDemand(name)
	if err != nil {
		return err
	}
	crawl := &Crawl{
		source: source,
		temp:   storage,
	}
	r.HandleFunc(fmt.Sprintf("/crawl/%s", source.Name()), crawl.FetchData).Methods("POST")
	return nil
}

var responseLogger = log.Sample(zerolog.LevelSampler{
	DebugSampler: &zerolog.BurstSampler{
		Burst:       5,
		Period:      1 * time.Second,
		NextSampler: &zerolog.BasicSampler{N: 100},
	},
})
