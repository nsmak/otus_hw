FROM golang:1.15.2 as builder

ENV BIN_FILE /opt/scheduler/scheduler-app
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
        -o ${BIN_FILE} cmd/scheduler/*

# На выходе тонкий образ
FROM alpine:3.9

ENV BIN_FILE "/opt/scheduler/scheduler-app"
COPY --from=builder ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/scheduler/config.json
COPY ./configs/scheduler.json ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}