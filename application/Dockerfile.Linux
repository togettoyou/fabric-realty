FROM golang:1.14
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
COPY . /root/application
WORKDIR /root/application
RUN go build -o "app" .