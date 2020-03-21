package storage

import (
	"context"
	"testing"
	"time"
)

func TestS3TempStorage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()

	db, err := TempDB(ctx)
	if err != nil {
		t.Fatal(err)
	}
	crawlTime := time.Now().UTC()
	expectedName, _ := ComputeObjectName("ministerio_saude_brasil", crawlTime, "json")
	name, err := db.StoreCrawl(ctx, "ministerio_saude_brasil", crawlTime, "json", []byte(`stuff`))
	if err != nil {
		t.Fatal(err)
	}
	if name != expectedName {
		t.Errorf("expecting name to be %v but got %v", expectedName, name)
	}
}

func TestComputeName(t *testing.T) {
	crawlTime := time.Unix(1577836800, 0).UTC().Round(time.Hour)
	t.Logf("crawlTime: %v", crawlTime.Format(time.RFC3339Nano))
	expectedName := "ministerio_saude_brasil/2020-01-01/00-00/rawData.json"
	actualName, err := ComputeObjectName("ministerio_saude_brasil", crawlTime, "json")
	if err != nil {
		t.Fatal(err)
	}
	if actualName != expectedName {
		t.Errorf("expecting actual name to be %v but got %v", expectedName, actualName)
	}
}
