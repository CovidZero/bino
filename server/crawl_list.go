package server

import (
	"net/http"

	"github.com/CovidZero/bino/datasources"
)

type (
	// CrawlerList retorna a lista de crawlers cadastrada atualmente
	CrawlerList struct {
	}
)

// GetAvailable retorna a lista de crawlers que está atualmente disponíveis
func (c CrawlerList) GetAvailable(w http.ResponseWriter, req *http.Request) {
	type (
		Definition struct {
			Name        string `json:"name"`
			Format      string `json:"format"`
			ContentType string `json:"contentType"`
			Encoding    string `json:"encoding"`
			Available   bool   `json:"available"`
		}

		Items struct {
			Crawlers []Definition `json:"crawlers"`
		}
	)

	var response Items
	for _, crawlName := range datasources.AllOnDemand() {
		ds, err := datasources.GetOnDemand(crawlName)
		def := Definition{
			Name: crawlName,
		}
		if err != nil {
			response.Crawlers = append(response.Crawlers, def)
			continue
		}
		def.Format = ds.Format().String()
		def.ContentType = ds.Format().ContentType()
		def.Encoding = ds.Encoding()
		def.Available = true

		response.Crawlers = append(response.Crawlers, def)
	}

	respondWithJSON(w, http.StatusOK, req, response)
}
