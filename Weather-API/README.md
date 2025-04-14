# Desafios FullCycle


## WEATHER API
Uma API escrita em Go que recebe um CEP, identifica a cidade correspondente e retorna o clima atual — com temperaturas em **Celsius**, **Fahrenheit** e **Kelvin**.

> Esta aplicação está publicada no **Google Cloud Run** e também pode ser executada via **Docker**.

---

### Como executar a aplicação com Docker

### ✅ Pré-requisitos

- Ter uma conta na [WeatherAPI](https://www.weatherapi.com/) e gerar sua **chave da API**.
- Ter o **Docker** instalado.

---

### Observação sobre o CEP

O parâmetro `cep` pode ser informado com ou sem traço:
- Sem o traço: `01001000`
- Com o traço: `01001-000`

---

### Passo a passo

1. **Clone o repositório**
```bash
git clone https://github.com/bianavic/fullcycle_desafios
```

2. Acesse o diretório do projeto
```bash
cd Weather-API
```

3. Construa a imagem Docker
```shell
docker build -t weather-api .
```

- Execute a aplicação com sua chave da WeatherAPI
  A aplicação não depende de um arquivo .env. A variável WEATHER_API_KEY pode (e deve) ser passada diretamente na linha de comando.

```shell
docker run -p 8080:8080 -e WEATHER_API_KEY=your_key weather-api
```

### Como acessar a aplicação no Cloud Run

Você também pode testar a aplicação publicada no Google Cloud Run:
```
https://weather-api-181365624128.us-central1.run.app/weather?cep=01001000
```
Substitua o cep no final da url pelo cep desejado

### Testes de Integração com .http
O projeto inclui um arquivo de testes prontos (weather-api-test.http) com requisições para Docker local e Cloud Run, 
ideal para quem usa extensões como o REST Client no VSCode. 
No IntelliJ IDEA (ou GoLand, da JetBrains) as requisições HTTP possuem suporte nativo.

Caminho do arquivo:
```bash
Weather-API/weather-api-test.http
```

### Exemplos de testes prontos:
✅ Sucesso com CEP válido

❌ Erro com CEP inválido

❌ Erro com CEP não informado