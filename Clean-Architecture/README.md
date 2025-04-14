# Desafio FullCycle

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
│ └── create_order.http
├── assets
│ └── images
│     ├── test1.png
│     ├── test2.png
│     ├── test3.png
│     ├── test4.png
│     └── test5.png
├── cmd
│ └── ordersystem
│     ├── main.go
│     ├── wire.go
│     └── wire_gen.go
├── configs
│ └── config.go
├── docker-compose.yaml
├── go.mod
├── go.sum
├── gqlgen.yml
├── internal
│ ├── entity
│ │ ├── interface.go
│ │ ├── order.go
│ │ └── order_test.go
│ ├── event
│ │ ├── handler
│ │ │ └── order_created_handler
│ │ │     └── order_created_handler.go
│ │ └── order_created.go
│ ├── infra
│ │ ├── database
│ │ │ ├── order_repository.go
│ │ │ ├── order_repository_test.go
│ │ │ └── sql
│ │ │     └── migrations
│ │ │         ├── 000001_init.down.sql
│ │ │         └── 000001_init.up.sql
│ │ ├── graph
│ │ │ ├── generated.go
│ │ │ ├── model
│ │ │ │ └── models_gen.go
│ │ │ ├── resolver.go
│ │ │ ├── schema.graphqls
│ │ │ └── schema.resolvers.go
│ │ ├── grpc
│ │ │ ├── pb
│ │ │ │ ├── order.pb.go
│ │ │ │ └── order_grpc.pb.go
│ │ │ ├── protofiles
│ │ │ │ └── order.proto
│ │ │ └── service
│ │ │     └── order_service.go
│ │ └── web
│ │     ├── order_handler.go
│ │     └── webserver
│ │         └── webserver.go
│ └── usecase
│     └── create_order.go
├── pkg
│ └── events
│     ├── event_dispatcher.go
│     ├── event_dispatcher_test.go
│     └── interface.go
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


- <img alt="imagem com o campo orders configurado no rabbit pq" src="assets/images/rabbitMQ1.png" width="434" height="356" title="Fila"/>

5b. Realizar bidding
- clicar na queue criada `orders` para realizar o bind à fila

- no campo `From exchange`, digitar `amq.direct` e clicar em Bind


- <img alt="imagem do bidding na aba exchanges" src="assets/images/rabbitMQ2.png" width="434" height="356" title="Fila"/>


### Request via webserver

- Create order
```shell
curl -X POST http://localhost:8000/order/create \
     -H "Content-Type: application/json" \
     -d '{
           "id": "aaaaaa",
           "price": 99.5,
           "tax": 0.5
         }'
```

`Order created: {aaaaaa 99.5 0.5 100}`

- List all orders
```shell
curl -X GET http://localhost:8000/order
```

`{"orders":[{"id":"bbbbbb","price":100,"tax":1,"final_price":101,"created_at":"2025-03-23 23:54:44 -03:00"}]}`

#### [rabbitMQ] Visualizar a mensagem gerada: 
- acessa: http://localhost:15672/#/queues/%2F/orders
- clicar em `Get messages` e verificar a mensagem

### Request via graphql

1. acessa porta 8080
```
http://localhost:8080/
```

2. cria mutation e executa
```shell
mutation createOrder {
  createOrder(input: {id: "bbbbbb", Price: 100, Tax: 1.0 }) {
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
      "Price": 100,
      "Tax": 1,
      "FinalPrice": 101
    }
  }
}
```

`Order created: {bbbbbb 100 1 101}`

### Request via grpc

1. No terminal, executar
```bash
evans -r repl
```

2. No prompt do evans, executar
```shell
package pb
usecase OrderService
call CreateOrder
```
responder no prompt grpc
```shell
id (TYPE_STRING) => cccccc
price (TYPE_FLOAT) => 12.2
tax (TYPE_FLOAT) => 2
```

```shell
{
  "finalPrice": 14.2,
  "id": "cccccc",
  "price": 12.2,
  "tax": 2
}
```

`Order created: {cccccc 12.199999809265137 2 14.199999809265137}`

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
