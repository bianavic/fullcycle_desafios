# Desafio FullCycle

## Objetivo
Criar um endpoint de listagem das orders para cada tipo de serviço:

- **Webserver**: Endpoint REST
- **GraphQL**: Query GraphQL
- **gRPC**: Endpoint gRPC

---

## Serviços e Portas

| Serviço     | Tipo      | Porta/URL                   |
|-------------|-----------|-----------------------------|
| Web Server  | REST      | `:8000`                     |
| gRPC        | Protobuf  | `:50051`                    |
| GraphQL     | Playground| [http://localhost:8081/](http://localhost:8081/) |

---

## Configuração e Execução

### Pré-requisitos
- Docker e Docker Compose instalados
- Go 1.23.8+ instalado
- Migrate CLI para execução de migrações

### Passo a Passo

1. **acessar raiz do projeto**
```shell
cd Clean-Architecture
```

2. subir containers docker
```shell
docker-compose up -d
```

2. rodar as migrações do banco de dados
```shell
migrate -path ./internal/infra/database/sql/migrations \
    -database "mysql://root:root@tcp(localhost:3306)/orders" \
    up
```

3. acessar o diretorio do serviço
```shell
cd cmd/ordersystem
```

4. iniciar servidores (Web, GraphQL, gRPC)
```shell
go run main.go wire_gen.go
```

---
### Request via webserver (REST)

- criar order
```shell
curl -X POST http://localhost:8000/order/create \
     -H "Content-Type: application/json" \
     -d '{
           "id": "eeeeee",
           "price": 20.0,
           "tax": 0.2
         }'
```

Resposta esperada:

`Order created: {aaaaaa 99.5 0.5 100}`

- listar orders
```shell
curl -X GET http://localhost:8000/order
```

Resposta esperada:

`{"orders":[{"id":"eeeeee","price":20,"tax":0.2,"final_price":20.2}]}`

---
### Request via gRPC

Recomendado usar o Evans como cliente CLI gRPC https://evans.syfm.me/

1. iniciar o cliente Evans
```bash
evans -r repl
```

2. no prompt do evans, executar um a um
```shell
package pb
service OrderService
```

Criar Order
```
call CreateOrder
```

Input
```shell
id (TYPE_STRING) => dddddd
price (TYPE_FLOAT) => 50
tax (TYPE_FLOAT) => 0.5
```

Resposta esperada:
```
{
  "finalPrice": 50.5,
  "id": "dddddd",
  "price": 50,
  "tax": 0.5
}
```

4. listar orders
```shell
call ListOrders
```

Resposta esperada:
```
{
  "orders": [
    {
      "finalPrice": 50.5,
      "id": "dddddd",
      "price": 50,
      "tax": 0.5
    },
    {
      "finalPrice": 20.2,
      "id": "eeeeee",
      "price": 20,
      "tax": 0.2
    }
  ]
}
```

---
### Request via GraphQL

1. acessar o playground
```
http://localhost:8081/
```

2. executar as chamadas

2a. criar order
```shell
mutation createOrder {
  createOrder(input: {id: "bbbbbb", price: 100, tax: 1.0 }) {
    id
    price
    tax
    finalPrice
  }
}
```

Resposta esperada:
```shell
{
  "data": {
    "createOrder": {
      "id": "bbbbbb",
      "Price": 100,
      "Tax": 1,
      "FinalPrice": 101
    }
  }
}
```

2b. listar orders

```shell
query {
  listOrders {
    id
    Price
    Tax
    FinalPrice
  }
}
```

Resposta esperada:
```shell
{
  "data": {
    "listOrders": [
      {
        "id": "bbbbbb",
        "Price": 99.9,
        "Tax": 0.6,
        "FinalPrice": 100.5
      },
      {
        "id": "dddddd",
        "Price": 50,
        "Tax": 0.5,
        "FinalPrice": 50.5
      },
      {
        "id": "eeeeee",
        "Price": 20,
        "Tax": 0.2,
        "FinalPrice": 20.2
      }
    ]
  }
}
```

---
### Notas Adicionais
- Certifique-se que todos os serviços estão rodando antes de fazer as requisições
- As migrações precisam ser executadas apenas na primeira execução
- Para reiniciar completamente, pare os containers com `docker-compose down -v` e recomece o processo
