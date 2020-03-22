package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CovidZero/bino/internal/testutil"
)

// TestIVISCrawl verifica se o processo para coletar os dados da base IVIS está funcionando como esperado
func TestIVISCrawl(t *testing.T) {
	ctx := context.TODO()

	fakeContent := "{}"
	ivisMockServer := testutil.ServeFile([]byte(fmt.Sprintf("var database=%v", fakeContent)))
	defer ivisMockServer.Close()
	defer testutil.ChangeEnv("DATASOURCE_IVIS_URL", ivisMockServer.URL)()

	db, getContent := testutil.TempDB(ctx, t)

	crawlers, err := allCollectors(db)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(crawlers)
	defer server.Close()
	response := struct {
		Path string `json:"path"`
	}{}
	status := testutil.POSTRaw(t, server.URL+"/crawl/ministerio_saude_brasil", "application/json", nil, &response, json.Unmarshal)
	if status != http.StatusOK {
		t.Fatalf("crawlers should always return 200 (at least for now...) but got %v", status)
	}
	t.Logf("Path: %v", response.Path)

	actualContent := string(getContent(ctx, t, response.Path))
	if actualContent != fakeContent {
		t.Fatalf("Should have %v in s3 but got %v", fakeContent, actualContent)
	}
}
