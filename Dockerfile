FROM golang:alpine AS builder


WORKDIR /src
COPY . .

RUN go build -o ./dist/grpc-rest ./cmd

FROM alpine:latest

WORKDIR /app
COPY --from=builder /src/dist/grpc-rest .

CMD ["./grpc-rest"]
