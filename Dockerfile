FROM dwarvesf/sql-migrate as sql-migrate

FROM golang:1.21-alpine as builder
RUN mkdir /build
WORKDIR /build

# manage app deps
COPY go.mod .
COPY go.sum .
RUN go mod download

# prepare base deps
ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0
RUN set -ex && \
    apk add --no-progress \
    musl-dev

# build
COPY . .
RUN go install -v --tags musl ./...

FROM alpine:3.18.0
RUN apk add ca-certificates
WORKDIR /
COPY --from=sql-migrate /usr/local/bin/sql-migrate /usr/bin/
COPY --from=builder /go/bin/* /usr/bin/
COPY migrations /migrations
COPY dbconfig.yml .

CMD [ "server" ]
