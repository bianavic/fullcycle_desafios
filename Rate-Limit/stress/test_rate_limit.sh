#!/bin/bash

# URL do servidor
URL="http://localhost:8082/"

# Token de acesso
TOKEN="your_token"

# Número de requisições a serem enviadas
NUM_REQUESTS=100

# Enviar múltiplas requisições
for i in $(seq 1 $NUM_REQUESTS); do
  curl -H "API_KEY: $TOKEN" $URL
  echo
done