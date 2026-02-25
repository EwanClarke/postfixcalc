FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -a -installsuffix cgo -o calc ./cmd/calc/ && \
    go build -a -installsuffix cgo -o plot ./cmd/plot/

FROM alpine:latest

RUN adduser -D -s /bin/sh appuser

WORKDIR /app

COPY --from=builder /app/calc .
COPY --from=builder /app/plot .
COPY --from=builder /app/entrypoint.sh .

RUN chown appuser:appuser calc plot entrypoint.sh
RUN chmod +x entrypoint.sh

USER appuser

CMD ["./entrypoint.sh"]
