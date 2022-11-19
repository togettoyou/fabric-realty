FROM golang:1.14 AS app
ENV GO111MODULE=on
ENV GOPROXY https://goproxy.cn,direct
WORKDIR /root/togettoyou
COPY server/. .
RUN CGO_ENABLED=0 go build -v -o "app" .

FROM node:16.17.1-alpine AS web
WORKDIR /root/togettoyou
COPY web/package*.json ./
RUN yarn config set registry https://registry.npmmirror.com/ && yarn install
COPY web/. .
RUN yarn run build:prod

FROM scratch
WORKDIR /root/togettoyou/
COPY --from=app /root/togettoyou/app ./
COPY --from=app /root/togettoyou/config.yaml ./
COPY --from=web /root/togettoyou/dist/ ./dist/
ENTRYPOINT ["./app"]