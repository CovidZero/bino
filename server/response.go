package server

import (
	"encoding/json"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, status int, req *http.Request, data interface{}) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	err := enc.Encode(data)
	if err != nil {
		// TODO incluir logging com um sample rate, para impedir DoS via Log, por hora, só ignora
		responseLogger.Debug().Msg("Error encoding data")
		return
	}

}
