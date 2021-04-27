FROM node:12.16.3-alpine
COPY . /root/vue
WORKDIR /root/vue
RUN yarn config set registry http://registry.npm.taobao.org \
  && yarn config set sass-binary-site http://npm.taobao.org/mirrors/node-sass \
  && yarn install
RUN yarn build:prod
