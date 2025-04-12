FROM golang:1.24-alpine AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun ./cmd/api

FROM scratch
WORKDIR /app
COPY --from=build /app/cloudrun .
ENTRYPOINT ["./cloudrun"]