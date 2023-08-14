FROM golang:1.20 as builder

RUN mkdir /app
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd


## [NOTE] with scratch, I can see below error
##       'There was an issue with sending an event: Post "https://xxxxxxxx.ingest.sentry.io/api/xxxxxxxx/envelope/": tls: failed to verify certificate: x509: certificate signed by unknown authority'
FROM alpine:3.18 as release
LABEL author="ghilbut@gmail.com"

EXPOSE 8080 9095

ENV GOMAXPROCS=1

COPY --from=builder /app/server /usr/local/bin/

CMD ["server"]
