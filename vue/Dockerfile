FROM node:12.16.3-alpine
COPY . /root/vue
WORKDIR /root/vue
RUN yarn config set registry https://registry.npmmirror.com \
  && yarn config set sass-binary-site http://npm.taobao.org/mirrors/node-sass \
  && yarn install
RUN yarn build:prod
