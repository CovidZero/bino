package main

import (
	"context"
	"flag"

	"github.com/CovidZero/bino/server"
	"github.com/CovidZero/bino/storage"
	"github.com/rs/zerolog/log"
)

var (
	bindAddr = flag.String("bind", ":8080", "Endereço onde a API irá rodar")
)

func main() {
	flag.Parse()

	// TODO: usar context.WithCancel e integrar com signal.Notify para fazer um shutdown limpo da aplicação
	rootCtx := context.Background()

	if errList := storage.SaneEnv(); errList != nil {
		for _, e := range errList {
			log.Error().Err(e).Msg("Check if your environment configuration is correct")
		}
		log.Fatal().Msg("Your environment is not properly configured. Please check the log for more information")
	}

	storage, err := storage.TempDB(rootCtx)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to start storage")
	}
	server, err := server.NewAPI(*bindAddr, storage)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to start server")
	}

	log.Info().Str("addr", *bindAddr).Msg("Starting server")
	// TODO: iniciar em um thread separado para ter suporte a tratar user signals
	if err := server.ListenAndServe(); err != nil {
		log.Error().Err(err).Msg("Server failed to exit")
	}
}
