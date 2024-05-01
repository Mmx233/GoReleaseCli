FROM golang:latest AS builder

WORKDIR /build

COPY . .

RUN go env -w GO111MODULE=auto \
    && go env -w CGO_ENABLED=0 \
    && set -ex \
    && go build -ldflags "-s -w" -o runner ./cmd/action

FROM alpine:latest

COPY --from=builder  /build/runner /usr/bin/runner
RUN chmod +x /usr/bin/runner

WORKDIR /data

ENTRYPOINT [ "/usr/bin/runner" ]