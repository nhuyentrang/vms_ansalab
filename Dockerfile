# syntax=docker/dockerfile:1

FROM --platform=${BUILDPLATFORM} golang:1.19-alpine3.15 AS builder

RUN apk add git

WORKDIR /go/src/app
COPY . .

ARG TARGETOS TARGETARCH TARGETVARIANT

ENV CGO_ENABLED=0
RUN go get \
    && go mod download \
    && GOOS=${TARGETOS} GOARCH=${TARGETARCH} GOARM=${TARGETVARIANT#"v"} go build -a -o vms_ansa

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /go/src/app/vms_ansa /app/
COPY --from=builder /go/src/app/web /app/web

RUN mkdir -p /config
COPY --from=builder /go/src/app/config.json /config

ENV GO111MODULE="on"
ENV GIN_MODE="release"

CMD ["./vms_ansa", "--config=/config/config.json"]
