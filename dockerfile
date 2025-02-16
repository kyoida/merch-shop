FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o merch-shop ./cmd/main.go

FROM scratch
COPY --from=builder /app/merch-shop /merch-shop
COPY --from=builder /app/.env /.env
EXPOSE 8080
ENTRYPOINT ["/merch-shop"]
