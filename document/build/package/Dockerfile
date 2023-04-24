FROM golang:1.20 as builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux
COPY go.* ./
RUN go mod download

COPY . .

RUN go build -a -installsuffix cgo -o documents document/cmd/app/main.go
RUN go build -a -installsuffix cgo -o migrate document/cmd/migrate/migrate.go

FROM alpine:latest

WORKDIR /build

ENV DOCUMENTS_APP_MIGRATIONS_DIR="migrations"

COPY --from=builder /app/documents /build/documents
COPY --from=builder /app/migrate /build/migrate
COPY --from=builder /app/document/internal/migrations/pgql /build/migrations

CMD ["./documents"]