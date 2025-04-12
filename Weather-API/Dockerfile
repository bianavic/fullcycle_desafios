FROM golang:1.24-alpine AS build

WORKDIR /app

# Instala certificados e outras libs necessárias
RUN apk update && apk add --no-cache ca-certificates

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun ./cmd/api

# Container final com certificados
FROM alpine:latest
# Copia os certificados de autoridade
RUN apk --no-cache add ca-certificates
#FROM scratch
WORKDIR /app

COPY --from=build /app/cloudrun .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["./cloudrun"]