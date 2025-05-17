FROM golang:1.24.3-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git ca-certificates gcc g++ && \
    update-ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/loadbalancer

RUN CGO_ENABLED=1 GOOS=linux go build -o ../../lb

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root

COPY --from=builder /app/lb .
COPY --from=builder /app/configs ./configs
ENTRYPOINT ["/root/lb"]