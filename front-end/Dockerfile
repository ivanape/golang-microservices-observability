FROM golang:bullseye as builder

RUN mkdir -p /app

COPY ./ /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o frontEnd ./cmd/api

RUN chmod +x /app/frontEnd

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/frontEnd /app

CMD ["/app/frontEnd"]