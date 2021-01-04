FROM golang:1.14
LABEL author sakura
LABEL version="1.0"
WORKDIR /go/src/app
ENV GOPROXY=https://proxy.golang.org,direct
COPY . .
RUN go mod download \
    &&go build -o crawl
CMD ./crawl

