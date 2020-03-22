package testutil

import (
	"net/http"
	"net/http/httptest"
)

// ServeFile retorna um server http que sempre retorna o conteudo ignorando o caminho
func ServeFile(content []byte) *httptest.Server {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(content)
	}))
	return s
}
