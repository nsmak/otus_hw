FROM golang:1.15.2
WORKDIR /app
COPY . .
RUN go test -i --tags=integration /app/internal/tests/integration/...
CMD go test -v --tags=integration /app/internal/tests/integration/...