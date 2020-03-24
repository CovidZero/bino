package testutil

import (
	"io/ioutil"
	"net/http"
	"testing"
)

// GET baixa a url e grava em out usando codec para decodificar
func GET(t *testing.T, url string, out interface{}, codec func([]byte, interface{}) error) int {
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Status %v / Error: %v", res.StatusCode, err)
	}
	println(string(data))
	err = codec(data, out)
	if err != nil {
		t.Fatalf("Status %v / Error: %v", res.StatusCode, err)
	}
	return res.StatusCode
}
