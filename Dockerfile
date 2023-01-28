FROM golang:1.19 as builder

ENV ENV=PROD
ENV GOPATH=/go
WORKDIR /build
COPY ./ ./

RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod download

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o web3_blog .

FROM debian:buster-slim

COPY --from=builder /build/web3_blog ./
COPY config/default.yaml ./config/
COPY templates ./templates

EXPOSE 9000

ENTRYPOINT ["./web3_blog"]