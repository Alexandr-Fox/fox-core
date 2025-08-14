FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# clear files
RUN make

FROM alpine:3 AS release

COPY --from=builder /app/build/* /usr/bin/fox-app

ENTRYPOINT ["/usr/bin/dumb-init", "/usr/bin/fox-app"]