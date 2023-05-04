FROM golang:1.20 as builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux
COPY go.* ./
RUN go mod download

COPY . .

RUN go build -a -installsuffix cgo -o sessions-app sessions/cmd/app/main.go

FROM alpine:latest

WORKDIR /build

COPY --from=builder /app/sessions-app /build/sessions
COPY --from=builder /app/api /build/api

CMD ["./sessions"]