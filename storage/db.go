package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"gocloud.dev/blob"
)

// TempDB retorna uma instância que pode ser usada para guardar dados temporários
// para processamento futuro
func TempDB(ctx context.Context) (Temp, error) {
	return newS3Storage(ctx)
}

func newS3Storage(ctx context.Context) (*S3Storage, error) {
	bucketName := os.Getenv("COVID0_TEMP_BUCKET")
	// INFO: https://gocloud.dev/howto/blob/#s3 contém informações de como a variável deve ser definida
	bucket, err := blob.OpenBucket(ctx, os.Getenv("COVID0_TEMP_BUCKET"))
	if err != nil {
		return nil, fmt.Errorf("unable to open bucket %s, cause: %v", bucketName, err)
	}
	return &S3Storage{
		bucket: bucket,
	}, nil
}

// SaneEnv verificar se as variáveis de ambiente exigidas pelo sistema estão devidamente configuradas
func SaneEnv() []error {
	// TODO: o uso de strings poderia ser trocado pelo uso de um bytes.Buffer, mas strings é mais fácil :-|
	var errList []error
	if bucket := os.Getenv("COVID0_TEMP_BUCKET"); bucket == "" {
		errList = append(errList, errors.New("missing COVID0_TEMP_BUCKET, please configure it so we know where to keep temporary files"))
	} else if !strings.Contains(bucket, "s3://") {
		errList = append(errList, errors.New("COVID0_TEMP_BUCKET is invalid, missing s3://. check https://gocloud.dev/howto/blob/#s3"))
	} else if !strings.Contains(bucket, "region=") {
		errList = append(errList, errors.New("COVID0_TEMP_BUCKET is invalid, missing region=<value>. check https://gocloud.dev/howto/blob/#s3"))
	}

	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		// INFO: checar se o caminho ~/.aws/credentials existe
		if path, err := homedir.Expand("~/.aws/credentials"); err != nil {
			errList = append(errList, fmt.Errorf("unable to open home dir, cause: %v", err))
		} else {
			_, err = os.Lstat(path)
			if os.IsNotExist(err) {
				errList = append(errList, fmt.Errorf(
					"missing AWS shared configuration at %v (check: https://docs.aws.amazon.com/sdk-for-go/api/aws/session/#hdr-Shared_Config_Fields), cause: %w",
					path, err))
			} else if err != nil {
				errList = append(errList, fmt.Errorf("unexpected error while checking path %v, cause %w", path, err))
			}
		}
	} else if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		errList = append(errList, errors.New("missing AWS environment variable AWS_SECRET_ACCESS_KEY"))
	}

	// TODO: verificar se as variáveis padrão da AWS estão configuradas
	return errList
}
