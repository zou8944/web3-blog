FROM golang:1.19 as builder

ENV ENV=PROD
ENV GOPATH=/go
WORKDIR /build
COPY ./ ./

RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o web3_blog .

FROM alpine

COPY --from=builder /build/web3_blog ./
COPY config/default.yaml ./config/
COPY templates ./templates

EXPOSE 9000

ENTRYPOINT ["./web3_blog"]