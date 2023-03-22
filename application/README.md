## 如改动代码需要自行编译后再使用 Docker 部署

后端：进入 `server` 执行 `./build.sh` ，并在 `docker-compose.yml` 配置本地镜像：`fabric-realty.server:latest`

前端：进入 `web` 执行 `./build.sh` ，并在 `docker-compose.yml` 配置本地镜像：`fabric-realty.web:latest`

## 支持本地开发模式

后端：更改 `server/blockchain/sdk.go` 中的配置文件路径为 `configPath = "config-local-dev.yaml"` 后，执行 `go run main.go`

前端：更改 `web/vue.config.js` 中的后端接口地址 `http://127.0.0.1:8888` 后，执行 `yarn install`
下载依赖，执行 `yarn run dev`