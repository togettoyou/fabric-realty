FROM node:12.16.3-alpine AS builder_vue
COPY vue /root/vue
WORKDIR /root/vue
RUN yarn config set registry http://registry.npm.taobao.org \
  && yarn config set sass-binary-site http://npm.taobao.org/mirrors/node-sass \
  && yarn install
RUN yarn build:prod

FROM golang:1.14
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
COPY application /root/application
COPY --from=builder_vue /root/vue/dist /root/application/dist
WORKDIR /root/application
RUN go build -o "app" .