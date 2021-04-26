FROM golang:1.14

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

RUN apt update \
    && apt-get -y install nodejs npm

COPY application /root/blockchain-real-estate/application

WORKDIR /root/blockchain-real-estate/application/vue
RUN npm config set registry "https://registry.npm.taobao.org/" \
  && npm i node-sass --sass_binary_site="https://npm.taobao.org/mirrors/node-sass/" \
  && npm install
RUN npm run build:prod

WORKDIR /root/blockchain-real-estate/application
RUN go build -o "app" .