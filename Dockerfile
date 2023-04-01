# Build Stage
FROM golang:1.19 AS build-stage

LABEL APP="build-grpc-gateway-example"
LABEL REPO="https://github.com/flamefatex/grpc-gateway-example"

ADD . /go/src/github.com/flamefatex/grpc-gateway-example
WORKDIR /go/src/github.com/flamefatex/grpc-gateway-example

RUN make build

# Final Stage
FROM alpine:3.17

ARG GIT_COMMIT
ARG VERSION
ARG APP_NAME

LABEL REPO="https://github.com/flamefatex/grpc-gateway-example"
LABEL GIT_COMMIT=${GIT_COMMIT}
LABEL VERSION=${VERSION}
LABEL APP_NAME=${APP_NAME}

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk add --no-cache tcpdump lsof net-tools tzdata curl dumb-init libc6-compat
RUN echo "hosts: files dns" > /etc/nsswitch.conf

ENV TZ Asia/Shanghai
ENV PATH=$PATH:/opt/grpc-gateway-example/bin

WORKDIR /opt/grpc-gateway-example/bin

COPY --from=build-stage /go/src/github.com/flamefatex/grpc-gateway-example/bin/grpc-gateway-example .
RUN chmod +x /opt/grpc-gateway-example/bin/grpc-gateway-example

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/grpc-gateway-example/bin/grpc-gateway-example"]
