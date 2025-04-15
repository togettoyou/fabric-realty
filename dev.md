# 本地开发指南

快速部署模式使用的是 Docker Hub 上已经编译好的镜像，适合快速体验系统功能。但在开发过程中，我们需要能够实时修改代码并查看效果，因此需要采用本地开发模式。

## 本地开发环境要求

- Go 1.23+
- Node.js 18+
- npm 9+
- Docker
- Docker Compose

## 本地开发步骤

### 1. 拉取项目（或手动下载）

```bash
git clone --depth 1 https://github.com/togettoyou/fabric-realty.git
```

### 2. 设置脚本权限

```bash
cd fabric-realty
find . -name "*.sh" -exec chmod +x {} \;
```

### 3. 启动区块链网络

首先需要启动基础的区块链网络环境（注意是进入到 network 目录执行）：

```bash
# 启动区块链网络（仅启动网络，不启动应用）
cd network
./install.sh
```

### 4. 启动后端服务

后端服务需要在本地编译运行，这样可以实时修改代码：

```bash
# 进入后端目录
cd application/server

# 运行后端服务
go run main.go
```

后端服务默认运行在 8888 端口。

### 5. 启动前端服务

前端服务同样需要在本地编译运行：

```bash
# 进入前端目录
cd application/web

# 安装依赖
npm install

# 运行开发服务器
npm run dev
```

前端开发服务器默认运行在 5173 端口。

### 6. 访问前端服务

http://localhost:5173

## 注意事项

1. 后端代码修改后，需要手动重启 `go run main.go`
2. 前端代码修改后，Vite 会自动热更新，无需手动重启
3. 区块链网络的修改（如链码更新）需要重新部署区块链网络
