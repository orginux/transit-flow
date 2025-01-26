FROM golang:1.23-alpine3.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /transit-flow

FROM alpine:3.21
RUN adduser -D appuser
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /transit-flow .
USER appuser
ENTRYPOINT ["./transit-flow"]
