FROM golang:1.21.0 AS builder

WORKDIR /storage/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o storage-service cmd/main.go

FROM debian:bookworm AS runner

WORKDIR /usr/bin

COPY --from=builder /storage/app/storage-service .
COPY --from=builder /storage/app/.env/ .

EXPOSE 8082

ENTRYPOINT ["storage-service"]