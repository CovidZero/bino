package testutil

import (
	"context"
	"os"
	"testing"

	"github.com/CovidZero/bino/storage"
	"gocloud.dev/blob"
)

// TempDB abre uma conexão com o banco de dados temporário.
//
// Note que o sistema presume que estará rodando da mesma forma que a aplicação roda,
// isso implica que você deve fornecer um servidor compatível com S3 e configurar a variável
// de ambiente COVID0_TEMP_BUCKET
func TempDB(ctx context.Context, t *testing.T) (storage.Temp, func(context.Context, *testing.T, string) []byte) {
	db, err := storage.TempDB(ctx)
	if err != nil {
		t.Fatal(err)
	}
	bucket, err := blob.OpenBucket(ctx, os.Getenv("COVID0_TEMP_BUCKET"))
	if err != nil {
		t.Fatal(err)
	}
	getfn := func(ctx context.Context, t *testing.T, key string) []byte {
		value, err := bucket.ReadAll(ctx, key)
		if err != nil {
			t.Fatal(err)
		}
		return value
	}
	return db, getfn
}
