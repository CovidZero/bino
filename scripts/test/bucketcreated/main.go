package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"gocloud.dev/blob"

	// suporte apenas para s3 por enquanto
	_ "gocloud.dev/blob/s3blob"
)

func main() {
	addr := flag.String("a", "", "bucket url")
	delay := flag.Duration("d", time.Second*10, "Time to wait between connection attempts")
	attempts := flag.Int("c", 10, "How many attempts before abort")
	flag.Parse()

	if len(*addr) == 0 {
		log.Fatal().Msg("Missing addr value")
	}

	if *attempts <= 0 {
		log.Error().Msg("Cannot use negative attempts, will default to 1")
		*attempts = 1
	}

	log.Info().Str("addr", *addr).Dur("delay", *delay).Int("attempts", *attempts).Send()

	sleep := func(err error) {
		log.Error().Err(err).Send()
		time.Sleep(*delay)
	}

	ctx := context.Background()
	for i := 0; i < *attempts; i++ {
		bucket, err := blob.OpenBucket(ctx, *addr)
		if err != nil {
			sleep(err)
			continue
		}
		err = createFile(ctx, bucket, "random_object")
		if err != nil {
			sleep(err)
			continue
		}
		bucket.Delete(ctx, "random_object")
		bucket.Close()
		log.Info().Str("addr", *addr).Dur("delay", *delay).Int("attempts", i+1).Msg("Success")
		return
	}
	os.Exit(2)
}

func createFile(ctx context.Context, bucket *blob.Bucket, key string) error {
	return bucket.WriteAll(ctx, key, []byte("hello"), &blob.WriterOptions{})
}
