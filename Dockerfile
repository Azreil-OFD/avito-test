FROM golang:1.22-alpine

WORKDIR ${GOPATH}/avito-shop/
COPY . ${GOPATH}/avito-shop/
RUN go mod tidy

RUN go build -o ./build ./cmd/service/ \
    && go clean -cache -modcache


CMD ["./build"]