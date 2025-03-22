# Desafios FullCycle


## Stress Test

Acesse o diretorio do projeto - Stress-Test

1. Construir imagem docker
```
docker build -t stress-test .
```

2.Executar aplicação
```
docker run stress-test --url=http://google.com --requests=1000 --concurrency=10
```

3.Executar testes
```shell
```

4. Resultado
```shell
```


### Configuração Cobra CLI

1. dependencias
```shell
go get -u github.com/spf13/cobra@latest
```
2. Inicialização
```shell
cobra-cli init <nome do projeto>
```
