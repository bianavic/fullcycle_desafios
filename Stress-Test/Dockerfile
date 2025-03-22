FROM golang:1.23-alpine

WORKDIR /app
COPY . .

# Build the Go app
RUN go build -o stress-test main.go

# Command to run the executable
ENTRYPOINT ["./stress-test"]