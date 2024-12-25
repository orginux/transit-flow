FROM golang:1.23-alpine3.21 AS build-stage
COPY go.mod go.sum ./
RUN go mod download
RUN go clean -cache -modcache -i -r
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /transit-flow

FROM alpine:3.21
WORKDIR /
COPY --from=build-stage /transit-flow /transit-flow
ENTRYPOINT ["/transit-flow"]
