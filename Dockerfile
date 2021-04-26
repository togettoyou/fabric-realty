FROM golang:1.14

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

RUN mkdir -p /root/blockchain-real-estate
RUN apt update \
    && apt-get -y install nodejs npm

COPY application /root/blockchain-real-estate/application
COPY deploy/crypto-config /root/blockchain-real-estate/crypto-config

WORKDIR /root/blockchain-real-estate/application/vue
RUN npm config set registry "https://registry.npm.taobao.org/" \
  && npm i node-sass --sass_binary_site="https://npm.taobao.org/mirrors/node-sass/" \
  && npm install
RUN npm run build:prod
RUN mv ./dist ../

WORKDIR /root/blockchain-real-estate/application
RUN go build -o "app" .

EXPOSE 8000
ENTRYPOINT ["./app"]