FROM golang:1.15.2 as builder
RUN go get -u github.com/pressly/goose/cmd/goose

FROM debian
WORKDIR /app
COPY --from=builder /go/bin/goose .
COPY ./migrations .
CMD sleep 10; ./goose postgres "postgres://postgres:password@db:5432/postgres?sslmode=disable" up;