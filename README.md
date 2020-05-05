# PROJETO ARQUIVADO E FORA DE UTILIZAÇÃO

# bino - transferindo dados de um lugar para outro

Bino é uma API que permite a seus usuários transferir dados de diferentes origens para um bucket S3.

O projeto em questão consiste em ler os dados relacionados ao COVID-19 do Ministério da Saude e do portal [Brasil.IO](https://brasil.io/home/) e inseri-los em um bucket S3

## API

```
curl -X POST http://localhost:8080/crawl/ministerio_saude_brasil
```

A operação acima irá baixar os dados do Ministério da Saúde e guardar em um bucket S3. A resposta segue o formato abaixo:

```
{ "path": "ministerio_saude_brasil/2006-01-02/15-04/rawData.json"}
```

```
curl -X GET http://localhost:8080/crawlers
```

A operação acima irá retornar uma lista de quais coletores estão configurados e disponíveis para execução

```js
{
	"crawlers": [{
		"name": "ministerio_saude_brasil",
		"format": "json",
		"contentType": "application/json",
		"encoding": "iso-8859-1",
		"available": true
	}, {
		"name": "brasil_io_covid19",
		"format": "csv",
		"contentType": "text/csv",
		"encoding": "utf-8",
		"available": true
	}]
}
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

Para rodar local, você pode usar o atalho `make run`

## Criando um novo datasource - OnDemand

Datasources `on-demand` são acionados via chamada URL e devem obter os dados a cada vez que são executados. Lembre-se
que por enquanto não existe proteção para impedir que ocorra exaustão de recursos externos. Sendo assim, verifique
a frequência que o endpoint será chamado e se é compatível com a capacidade do datasource que você irá incluir.

O formato dos dados não é importante, `bino` opera com dados bytes e o processamento feito deve ser o mínimo possível,
normalmente remover alguns sufixos ou prefixos.

Com isso em mente, vamos aos passos que devem ser feitos para incluir um novo Datasource HTTP.

- Na pasta `datasources` crie um arquivo `nome_do_meu_datasource.go` (tente usar nomes curtos).
- No arquivo `nome_do_meu_datasource.go` escreva o seu coletor (use `ivis.go` como exemplo)
- No arquivo `ondemand.go` altere as funções para `AllOnDemand` e `GetOnDemand` para expor o seu datasource
- Se os dados do seu datasource não são JSON, adicione o formato novo em `format.go`
- Na pasta `server` escreva um teste de integração para o seu datasource (use `crawl_ivis_test.go` como exemplo)
- Evite criar muitos mocks, Go é uma linguagem compilada e mocks em encesso não são a melhor maneira de trabalhar.
- Descreva em [Datasources.md](Datasources.md) o seu datasource.
- Be Happy e fique em casa! Sério... Fique em casa! :-)
