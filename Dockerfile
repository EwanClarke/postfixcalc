FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -a -installsuffix cgo -o postfixcalc ./cmd/calc/

FROM alpine:latest

RUN adduser -D -s /bin/sh appuser

WORKDIR /app

COPY --from=builder /app/postfixcalc .

RUN chown appuser:appuser postfixcalc

USER appuser

ENTRYPOINT ["./postfixcalc"]
