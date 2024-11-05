# Use the official Go image as the base image
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o vault-cache .

# Use a smaller base image for the final stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/vault-cache .

ENV upstream=http://localhost:8200
ENV port=8201
ENV ttl=5m

EXPOSE $port

CMD /app/vault-cache -upstream ${upstream} -port ${port} -ttl ${ttl}
