package testutil

import (
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

// POSTRaw executa uma chamada HTTP post na URL informada com o body e contentType informados.
//
// Caso uma resposta seja obtida, o conteúdo é lido para memória e depois enviado para a função
// codec que deve transformar os bytes recebidos e guardar em out
func POSTRaw(t *testing.T, url string, contentType string, body io.Reader, out interface{}, codec func([]byte, interface{}) error) int {
	res, err := http.Post(url, contentType, body)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	err = codec(data, out)
	if err != nil {
		t.Fatal(err)
	}
	return res.StatusCode
}
