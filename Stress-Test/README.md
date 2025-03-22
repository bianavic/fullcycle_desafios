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
```shell
docker run stress-test --url=http://test.com --requests=1000 --concurrency=10
```

3. Execute o comando abaixo para executar o teste unitario
```shell
./stress-test --url=http://test.com --requests=1000 --concurrency=10
```

4. Relatorio de teste
```shell
```
