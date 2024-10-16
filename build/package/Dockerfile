FROM golang:1.22 as builder

WORKDIR $GOPATH/src/dns-server
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /dns-server -x cmd/main.go

FROM alpine:3.19

COPY --from=0 /dns-server /usr/bin
WORKDIR /

ENTRYPOINT ["sh", "-c", "dns-server"]
