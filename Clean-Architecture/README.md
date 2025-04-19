# Desafio FullCycle

obs: erro para grapql createorder

## Listagem das ordens - criar um enpoint para cada tipo de serviço
- via webserver - criar endpoint rest
- via graphql - criar query graphql
- via grpc - criar endpoint grpc

Project Structure
```
.
├── LICENSE
├── README.md
├── api
│   ├── create_order.http
│   └── list_orders.http
├── assets
│   └── images
│       ├── rabbitMQ1.png
│       ├── rabbitMQ2.png
│       ├── rabbitMQ3.png
│       ├── test1.png
│       ├── test2.png
│       ├── test3.png
│       ├── test4.png
│       └── test5.png
├── cmd
│   └── ordersystem
│       ├── main.go
│       ├── wire.go
│       └── wire_gen.go
├── configs
│   └── config.go
├── docker-compose.yaml
├── go.mod
├── go.sum
├── internal
│   ├── dto
│   │   └── order_dto.go
│   ├── entity
│   │   ├── interfaces.go
│   │   ├── order.go
│   │   └── order_test.go
│   ├── events
│   │   ├── handlers
│   │   │   └── order_created_handler.go
│   │   ├── order_created.go
│   │   └── order_listed.go
│   ├── infrastructure
│   │   ├── database
│   │   │   ├── order_repository.go
│   │   │   ├── order_repository_test.go
│   │   │   └── migrations
│   │   │       ├── 000001_init.down.sql
│   │   │       └── 000001_init.up.sql
│   │   ├── graphql
│   │   │   ├── generated.go
│   │   │   ├── model
│   │   │   │   └── models_gen.go
│   │   │   ├── resolver.go
│   │   │   ├── schema.graphqls
│   │   │   └── schema.resolvers.go
│   │   ├── grpc
│   │   │   ├── pb
│   │   │   │   ├── order.pb.go
│   │   │   │   └── order_grpc.pb.go
│   │   │   ├── proto
│   │   │   │   └── order.proto
│   │   │   └── service
│   │   │       └── order_service.go
│   │   └── http
│   │       ├── handlers
│   │       │   └── order_handler.go
│   │       └── server.go
│   └── application
│       ├── commands
│       │   └── create_order.go
│       └── queries
│           └── list_orders.go
├── pkg
│   └── events
│       ├── dispatcher.go
│       ├── dispatcher_test.go
│       └── interfaces.go
└── tools.go
```

## Executar

1. subir docker
```shell
docker-compose up -d
```

2. Rodar migração
```
migrate -path ./internal/infra/database/sql/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" up
```

3. acessa diretorio
```shell
cd cmd/ordersystem
```

4. executa servidores (webserver, graphql, grpc)
```shell
go run main.go wire_gen.go
```

***
[NOTA] Configure rabbitMQ antes de executar as requests
***

5. Preparar rabbitMQ

5a. Acesso
```
http://localhost:15672/
```
senha: guest

5b. Seleciona a aba Queues and Streams
```
http://localhost:15672/#/queues
```

5c. Criar fila
- no campo `name`, digitar orders e clicar em Add queue

<img alt="imagem com o campo orders configurado no rabbit pq" src="assets/images/rabbitMQ1.png" width="434" height="356" title="Fila"/>

5b. Realizar bidding
- clicar na queue criada `orders` para realizar o bind à fila
- no campo `From exchange`, digitar `amq.direct` e clicar em Bind

<img alt="imagem do bidding na aba exchanges" src="assets/images/rabbitMQ2.png" width="434" height="356" title="Fila"/>

6b. Verificar a mensagem
***Realizar chamada apenas após criar Order***
- Aba `Queues and Streams`
- clicar em `Get messages` e verificar a mensagem

<img alt="imagem do rabbitMQ com a mensagem" src="assets/images/rabbitMQ3.png" width="434" height="356" title="Fila"/>


### Request via webserver

- Create order
```shell
curl -X POST http://localhost:8000/order/create \
     -H "Content-Type: application/json" \
     -d '{
           "id": "eeeeee",
           "price": 20.0,
           "tax": 0.2
         }'
```

`Order created: {aaaaaa 99.5 0.5 100}`

- List all orders
```shell
curl -X GET http://localhost:8000/order
```

`{"orders":[{"id":"eeeeee","price":20,"tax":0.2,"final_price":20.2}]}`

### Request via grpc

1. No terminal, executar
```bash
evans -r repl
```

2. No prompt do evans, executar
```shell
package pb
```

3. Acessar service, executar
````shell
service OrderService
````

3. Criar Order, executar
````shell
call CreateOrder
````

4. Responder no prompt grpc
```shell
id (TYPE_STRING) => dddddd
price (TYPE_FLOAT) => 50
tax (TYPE_FLOAT) => 0.5

```

Response esperada
```shell
{
  "finalPrice": 50.5,
  "id": "dddddd",
  "price": 50,
  "tax": 0.5
}
```

5. Listar orders, executar
```shell
call ListOrders
```

Response esperada
```shell
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

### Request via graphql

1. acessa porta 8080
```
http://localhost:8081/
```

2. cria Order e executa
```shell
mutation {
  createOrder(input: {
    id: "bbbbbb"
    Price: 99.9
    Tax: 0.6
  }) {
    id
    Price
    Tax
    FinalPrice
  }
}
```

output
```shell
{
  "data": {
    "createOrder": {
      "id": "bbbbbb",
      "Price": 99.9,
      "Tax": 0.6,
      "FinalPrice": 100.5
    }
  }
}
```

3. Lista orders e executa
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

output
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

#### [rabbitMQ] Visualizar a mensagem gerada: 
- acessa: http://localhost:15672/#/queues/%2F/orders
- clicar em `Get messages` e verificar a mensagem

### Request via graphql

1. acessa porta 8081
```
http://localhost:8081/
```

2. cria mutation e executa
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

output
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

3. Lista orders e executa
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

output
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

## Banco MySQL - visualizar dados
1. acessar banco mysql
```
docker exec -it mysql bash
```

2. acessar banco orders
```
mysql -u root -p orders
```
senha: root

3. Ver tabela
```
select * from orders;
```
---
