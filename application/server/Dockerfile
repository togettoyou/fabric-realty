FROM golang:1.18 AS server
ENV GO111MODULE=on
ENV GOPROXY https://goproxy.cn,direct
WORKDIR /root/togettoyou
COPY . .
RUN CGO_ENABLED=0 go build -v -o "server" .

FROM scratch
WORKDIR /root/togettoyou/
COPY --from=server /root/togettoyou/server ./
COPY config.yaml ./
ENTRYPOINT ["./server"]