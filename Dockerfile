FROM golang:1.10.0 AS build-env

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src/github.com/lox/parfait
ADD . /go/src/github.com/lox/parfait
RUN go build -a -tags netgo -ldflags '-w' -o /bin/parfait

FROM scratch
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-env /bin/parfait /parfait
ENTRYPOINT ["/parfait"]
