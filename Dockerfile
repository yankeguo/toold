ARG VERSION="dev"

FROM golang:1.23 AS builder

ARG VERSION

ENV CGO_ENABLED="0"

WORKDIR /go/src/app

ADD . .

RUN go build -ldflags="-X main.AppVersion=${VERSION}" -o /toold ./cmd/toold

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tini

COPY --from=builder /toold /usr/local/bin/toold

ENTRYPOINT [ "/sbin/tini", "--" ]

CMD ["/usr/local/bin/toold"]

WORKDIR /toold