# Desafio FullCycle

Neste desafio você terá que usar o que aprendemos com Multithreading e APIs para buscar o resultado mais rápido entre duas APIs distintas.

As duas requisições serão feitas simultaneamente para as seguintes APIs:

https://brasilapi.com.br/api/cep/v1/01153000 + cep

http://viacep.com.br/ws/" + cep + "/json/

## Requisitos:

- Acatar a API que entregar a resposta mais rápida e descartar a resposta mais lenta.

- O resultado da request deverá ser exibido no command line com os dados do endereço, bem como qual API a enviou.

- Limitar o tempo de resposta em 1 segundo. Caso contrário, o erro de timeout deve ser exibido.

## How To

1- clone o respositorio
```
git clone https://github.com/bianavic/fullcycle_desafios.git
```

2- acesse a pasta do desafio
```shell
cd fullcycle_desafios/Multithreading
```

3- instale as dependencias
```shell
go mod tidy
```

4- rode o comando com o cep desejado (exemplo abaixo)
```shell
go run main.go 01153000
```

#### Expected:

Para simular API BrasilAPI com resposta mais rápida, configure o *tempo de resposta da thread da ViaCep para:
- ser maior que a thread da BrasilAPI
- ser menor que o valor de apiTimeout
```go
if source == "ViaCepAPI" {
    time.Sleep(time.Second * 4) // Simulate delay for ViaCepAPI
}
```
output:
```go
Received from brasilapi: source:BrasilAPI - Rua Vitorino Carmilo, Barra Funda - São Paulo, SP, 01153000
```

Para simular API ViaCep com resposta mais rápida, configure o *tempo de resposta da thread BrasilAPI para:
- ser maior que o da thread da ViaCep.
- Ser menor que o valor de apiTimeout.
```go
if source == "BrasilAPI" {
 time.Sleep(time.Second * 4) // Simulate delay for BrasilAPI
}
```
  output:
```go
Received from viacep: source:ViaCepAPI - Rua Vitorino Carmilo, Barra Funda - São Paulo, São Paulo, 01153-000
```

Para simular o tempo de resposta excedendo o limite, configure o *tempo de resposta para ambas as threads:
- maior que o valor de apiTimeout.
```go
if source == "BrasilAPI" {
 time.Sleep(time.Second * 4) // Simulate delay for BrasilAPI
 address = &BrasilAPIResponse{}
 } else {
 time.Sleep(time.Second * 4) // Simulate delay for ViaCepAPI
 address = &ViaCepAPIResponse{}
}
```
output:
```go
Timeout
```

### Tools
