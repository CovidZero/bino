package server

import (
	"net/http"
	"time"

	"github.com/andrebq/covid0-backend/storage"
	"github.com/gorilla/mux"
)

// NewAPI retorna um servidor HTTP pré-configurado que guarda as informações temporárias no storage informado
func NewAPI(bindAddr string, storage storage.Temp) (*http.Server, error) {
	crawl := &Crawl{temp: storage}

	r := mux.NewRouter()
	r.HandleFunc("/crawl/{source}", crawl.FetchData).Methods("POST")

	server := &http.Server{
		Addr: bindAddr,
		// Proteger o server de clients muito lentos, evitando que existam muitas conexões TCP abertas ao mesmo tempo
		// IMPORTANTE! não protege de clientes que intencionamente enviam dados a uma velocidade muito lenta (Slow Loris Attack)
		// https://www.youtube.com/watch?v=XiFkyR35v2Y para mais informações
		ReadTimeout: time.Second * 10,

		// Mesma lógica acima e mesma limitação
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 * 1024 * 1024 * 1024,

		Handler: r,
	}
	return server, nil
}
