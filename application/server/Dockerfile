FROM golang:1.23-alpine AS builder
ENV GOPROXY https://goproxy.cn,direct
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -mod=vendor -o main

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/config/config-docker.yaml /app/config/config.yaml
CMD ["./main"]
