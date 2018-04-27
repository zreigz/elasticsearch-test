FROM golang:1.7.3 as builder
WORKDIR /go/src/github.com/zreigz/elasticsearch-test/
RUN go get gopkg.in/olivere/elastic.v3 github.com/Pallinder/go-randomdata
COPY elastictest.go .
RUN ls -al
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o elasticsearch-test .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/zreigz/elasticsearch-test/elasticsearch-test .
ENTRYPOINT ["./elasticsearch-test"]
CMD ["--server", "http://127.0.0.1:9200"]
