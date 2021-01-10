FROM golang:1.15.2 as builder

ENV BIN_FILE /opt/sender/sender-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/sender/*

# На выходе тонкий образ
FROM alpine:3.9

ENV BIN_FILE "/opt/sender/sender-app"
COPY --from=builder ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/sender/config.json
COPY ./configs/sender.json ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}