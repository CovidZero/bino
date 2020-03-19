package server

import (
	"net/http"

	"github.com/andrebq/covid0-backend/storage"
)

type (
	// Crawl expões as operações de baixar e guardar arquivos em um storage temporário
	Crawl struct {
		temp storage.Temp
	}
)

// FetchData inicia um novo processo de coleta (caso exista necessidade), caso contrário,
// retorna a última coleta feita
func (c *Crawl) FetchData(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Ainda sem fazer nada", http.StatusNoContent)
}
