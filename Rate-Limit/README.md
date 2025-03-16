# Desafio FullCycle

### Objetivo
Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

```
.
├── README.md
├── cmd
│   └── api
│       └── main.go
├── docker-compose.yaml
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   └── config.go
│   ├── ratelimit
│   │   ├── middleware.go
│   │   └── middleware_test.go
│   ├── repository
│   │   └── storage
│   │       ├── in_memory.go
│   │       ├── mock
│   │       │   ├── mock_redis_client.go
│   │       │   └── mock_storage_strategy.go
│   │       ├── redis.go
│   │       ├── redis_client.go
│   │       ├── redis_test.go
│   │       └── storage_strategy.go
│   └── usecase
│       ├── rate_limiter.go
│       └── test
│           └── rate_limiter_test.go
└── stress
    ├── test_rate_limit.sh
    └── test_rate_limit_over.sh
```


## Executando a Aplicação

1. Clone o repositório.
2. Navegue até o diretório `Rate-Limit`
3. Execute ` docker-compose up --build` para iniciar os contêineres

## Testando

1. Limitação por IP:
```bash
for i in {1..6}; do curl -X GET http://localhost:8080/; done
```
A sexta requisição deve retornar 429 Too Many Requests. Aguarde 60 segundos para poder fazer outra requisição.

2. Limitação por Token:

- Token 1 (10 requisições por 60 segundos):
```bash
for i in {1..11}; do curl -X GET -H "API_KEY: abc123" http://localhost:8080/; done
```
A décima primeira requisição deve retornar 429 Too Many Requests. Aguarde 60 segundos para poder fazer outra requisição.

- Token 2 (20 requisições por 65 segundos):
```bash
for i in {1..21}; do curl -X GET -H "API_KEY: def456" http://localhost:8080/; done
```
A vigésima primeira requisição deve retornar 429 Too Many Requests. Aguarde 65 segundos para poder fazer outra requisição.

## Redis commander
1. Navegue até `http://127.0.0.1:8081/` 
2. A cada requisição, dê um refresh na página para visualizar detalhes relacionados ao rate limit.
![redis_commander1.png](assets/images/redis_commander1.png)

### Testes Automatizados:
localizados em `internal/usecase/test/rate_limiter_test.go`, `internal/repository/storage/redis_test.go` e `internal/middleware/middleware_test.go`