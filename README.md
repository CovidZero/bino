# bino - transferindo dados de um lugar para outro

Bino é uma API que permite a seus usuários transferir dados de diferentes origens para um bucket S3.

## API

```
curl -X POST http://localhost:8080/crawl/ministerio_saude_brasil
```

A operação acima irá baixar os dados do Ministério da Saúde e guardar em um bucket S3. A resposta segue o formato abaixo:

```
{ "path": "ministerio_saude_brasil/2006-01-02/15-04/rawData.json"}
```

## Setup do ambiente de dev

- Unix based (macos, wsl2, linux)
- Go SDK (1.13 no mínimo)
- Rode `make devsetup` (instala um git-hook para pre-commit)

## Compilando

Tenha a versão 1.13 da SDK GO e rode `make build`

## Testando

`make test` sobe o ambiente de suporte de teste (localstack) e roda todos os tests da aplicação.

Detalhes da execução estão em `scripts/test/integration_test.sh`

## Empacotando

Rode `make package`, você precisa ter o docker instalado e rodando

## Rodando a aplicação

Empacote a aplicação usando docker como descrito acima e execute o seguinte comando

```
# defina as variáveis AWS_ACCESS_KEY_ID e AWS_SECRET_ACCESS_KEY seguindo a recomendação abaixo
# https://docs.aws.amazon.com/sdk-for-go/api/aws/session/#hdr-Environment_Variables
#
# defina a variável COVID0_TEMP_BUCKET seguindo a recomendação abaixo
# https://gocloud.dev/howto/blob/#s3
# exemplo:
# COVID0_TEMP_BUCKET=s3://nome-do-bucket?region=<região do bucket>
#
# Para mais detalhes veja o arquivo storage/db.go função SaneEnv
docker run -p '8080:8080' -e COVID0_TEMP_BUCKET -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY covidzero/bino:latest
```

Para fins de teste, você pode usar o atalho `make run`
