## Builder stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o switchyard-go ./cmd/main.go


## Runtime stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/switchyard-go         ./switchyard-go
COPY --from=builder /build/internal/migrations   ./internal/migrations
COPY --from=builder /build/web/templates         ./web/templates
COPY --from=builder /build/web/static            ./web/static

RUN mkdir -p ./output/invoices

EXPOSE 8080

ENTRYPOINT ["./switchyard-go"]
