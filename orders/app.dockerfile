FROM golang:1.26-alpine3.23 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/ShivankSharma070/go-microservices-ecommerce

COPY go.mod go.sum ./
COPY vendor vendor
COPY account account
COPY catalog catalog
COPY orders orders

RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./orders/cmd/orders

FROM alpine:3.23
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]
