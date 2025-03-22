# Desafios FullCycle


## Stress Test

1. acesse o diretorio do projeto - Stress-Test
```shell
cd Stress-Test
```
2. Execute o comando abaixo para construir a imagem docker
```shell
docker build -t stress-test .
```
3. Execute o comando abaixo para executar o teste de stress com docker

- google.com é o endereço do site que será testado
```shell
docker run stress-test -u http://google.com -r 1000 -c 10
```

- uol.com.br é o endereço do site que será testado
```shell
docker run stress-test -u http://uol.com.br -r 100 -c 10
```

Exemplo execução e relatório gerado

![stress-test1.png](assets/images/stress-test1.png)