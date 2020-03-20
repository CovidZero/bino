package main

import (
	"flag"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {
	addr := flag.String("a", "", "TCP address to wait for")
	delay := flag.Duration("d", time.Second*10, "Time to wait between connection attempts")
	attempts := flag.Int("c", 10, "How many attempts before abort")
	flag.Parse()

	if len(*addr) == 0 {
		log.Fatal().Msg("Missing addr value")
	}

	if *attempts <= 0 {
		log.Error().Msg("Cannot use negative attempts, will default to 1")
		*attempts = 1
	}

	log.Info().Str("addr", *addr).Dur("delay", *delay).Int("attempts", *attempts).Send()

	for i := 0; i < *attempts; i++ {
		conn, err := net.Dial("tcp", *addr)
		if err != nil {
			log.Error().Err(err).Send()
			time.Sleep(*delay)
		} else {
			log.Info().Str("addr", *addr).Dur("delay", *delay).Int("attempts", i+1).Msg("Success")
			conn.Close()
			return
		}
	}
	os.Exit(2)
}
