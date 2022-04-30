FROM golang:1.17 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o jobs cmd/jobs/main.go

FROM alpine:latest AS production

COPY --from=builder /app .
CMD ["./jobs"]
