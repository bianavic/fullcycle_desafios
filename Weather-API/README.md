# Desafios FullCycle


## WEATHER API

URL: https://weather-api-181365624128.us-central1.run.app

1. Redeploy no cloud run com docker
```
gcloud builds submit --tag gcr.io/weather-api-456523/weather-api
gcloud run deploy weather-api \
  --image gcr.io/weather-api-456523/weather-api \
  --region=us-central1 \
  --platform=managed \
  --allow-unauthenticated

```

Teste weather api key
```
curl "http://api.weatherapi.com/v1/current.json?key=$WEATHER_API_KEY&q=São Paulo"
```