package main

import (
	"flag"

	"github.com/andrebq/covid0-backend/server"
	"github.com/andrebq/covid0-backend/storage"
	"github.com/rs/zerolog/log"
)

var (
	bindAddr = flag.String("bind", ":8080", "Endereço onde a API irá rodar")
)

func main() {
	flag.Parse()

	storage, err := storage.TempDB()
	if err != nil {
		// não faz sentido iniciar a aplicação se não foi possível abrir o banco temporário
		panic(err)
	}
	server, err := server.NewAPI(*bindAddr, storage)
	if err != nil {
		// não faz sentido iniciar a aplicação se não foi possível abrir o socket
		panic(err)
	}

	log.Info().Str("addr", *bindAddr).Msg("Starting server")
	// TODO: iniciar em um thread separado para ter suporte a tratar user signals
	if err := server.ListenAndServe(); err != nil {
		log.Error().Err(err).Msg("Server failed to exit")
	}
}
