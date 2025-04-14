FROM golang:1.24-alpine AS build

WORKDIR /app

RUN apk update && apk add --no-cache ca-certificates

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=build /app/cloudrun .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["./cloudrun"]