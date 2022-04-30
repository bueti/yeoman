FROM golang:1.17 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o jobs cmd/jobs/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

FROM alpine:latest AS job

COPY --from=builder /app .
CMD ["./jobs"]

FROM alpine:latest AS server

COPY --from=builder /app .
CMD ["./server"]
