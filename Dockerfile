FROM golang:1.20-alpine
RUN mkdir /build
WORKDIR /build
COPY . .
ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0
RUN set -ex && \
    apk add --no-progress --no-cache \
    gcc \
    musl-dev
RUN go install ./...
RUN go install github.com/rubenv/sql-migrate/sql-migrate@latest

FROM alpine:3.15.0
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=0 /go/bin/* /usr/bin/
COPY migrations /migrations
COPY dbconfig.yml .

CMD [ "server" ]
