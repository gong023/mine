FROM golang:1.17-alpine as builder

ARG CMD_PATH
ENV CGO_ENABLED 0

WORKDIR /opt/mine

ADD . /opt/mine

RUN go build -o /app ./cmd/server

FROM alpine

RUN apk add --update --no-cache ca-certificates

COPY --from=builder /app /app

ENTRYPOINT ["/app"]
