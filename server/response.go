package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

func respondWithJSON(w io.Writer, req *http.Request, data interface{}) {
	enc := json.NewEncoder(w)
	err := enc.Encode(data)
	if err != nil {
		// TODO incluir logging com um sample rate, para impedir DoS via Log, por hora, só ignora
		log.Error().Err(err).Str("module", "response").Msg("Error encoding data")
		return
	}

}
