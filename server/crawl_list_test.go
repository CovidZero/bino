package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CovidZero/bino/datasources"
	"github.com/CovidZero/bino/internal/testutil"
)

// TestCrawlerList verifica se todos crawlers esperados estão disponíveis
func TestCrawlerList(t *testing.T) {
	ctx := context.TODO()

	fakeContent := "{}"
	ivisMockServer := testutil.ServeFile([]byte(fmt.Sprintf("var database=%v", fakeContent)))
	defer ivisMockServer.Close()
	defer testutil.ChangeEnv("DATASOURCE_IVIS_URL", ivisMockServer.URL)()

	db, _ := testutil.TempDB(ctx, t)

	crawlers, err := allCollectors(db)
	if err != nil {
		t.Fatal(err)
	}

	var response struct {
		Crawlers []struct {
			Name        string `json:"name"`
			Format      string `json:"format"`
			ContentType string `json:"contentType"`
			Encoding    string `json:"encoding"`
			Available   bool   `json:"available"`
		} `json:"crawlers"`
	}
	server := httptest.NewServer(crawlers)
	defer server.Close()
	status := testutil.GET(t, server.URL+"/crawlers", &response, json.Unmarshal)
	if status != http.StatusOK {
		t.Fatalf("crawlers should always return 200 (at least for now...) but got %v", status)
	}
	t.Logf("Response: %v", response)
	expectedCrawlers := map[string]bool{
		datasources.IVISDatasetName:            true,
		datasources.BrasilIOCovid19DatasetName: true,
	}
	if len(response.Crawlers) != len(expectedCrawlers) {
		t.Errorf("Should have %v crawlers available got %v", len(expectedCrawlers), len(response.Crawlers))
	}
	for _, v := range response.Crawlers {
		expectedAvailability := expectedCrawlers[v.Name]
		if expectedAvailability != v.Available {
			t.Errorf("Crawler %v should have availability=%v but got %v", v.Name, expectedAvailability, v.Available)
		}
	}
}
