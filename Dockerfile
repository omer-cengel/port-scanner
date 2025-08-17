# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o port-scanner ./cmd/main.go

# Final stage
FROM scratch
WORKDIR /app

COPY --from=builder /app/port-scanner .

ENV DOCKERIZED=true

ENTRYPOINT ["./port-scanner"]
