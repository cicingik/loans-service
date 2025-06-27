FROM golang:1.22.4-alpine3.20 as builder
ARG COMMIT_MSG
ARG COMMIT
ARG GIT_TAG
ARG BUILD_TIME
ARG BUILD_TAG

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux

WORKDIR /app
COPY . /app

RUN set -e \
    && cd /app \
    && go mod download \
    && echo "Building binary..." \
        && go build -a -v -i \
            -ldflags "-w \
                -X \"github.com/cicingik/loans-service/config.BuildTime=${BUILD_TIME}\" \
                -X \"github.com/cicingik/loans-service/config.CommitMsg=${COMMIT_MSG}\" \
                -X \"github.com/cicingik/loans-service/config.CommitHash=${COMMIT}\" \
                -X \"github.com/cicingik/loans-service/config.AppVersion=${GIT_TAG}\" \
                -X \"github.com/cicingik/loans-service/config.ReleaseVersion=${BUILD_TAG}\"" \
            -tags ${BUILD_TAG} \
            -o loans

FROM alpine:3.20
LABEL author="dany <danysatyanegara@gmail.com>" \
      version="v1.1.0-beta" \
      description="docker image for loans-service"
ARG CO_HTTP_PORT=2777
ENV CO_HTTP_PORT_BIND=$CO_HTTP_PORT

WORKDIR /app/
COPY --from=builder /app/loans .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

EXPOSE ${CO_HTTP_PORT_BIND}
CMD [ "/app/loans" ]
