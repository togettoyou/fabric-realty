# 1. 安装 Docker CE（即 Docker 社区版）

使用curl下载并安装脚本

```shell
curl -sSL https://get.daocloud.io/docker | sh
```

（如果是root用户，请忽略）普通用户需要设置成非root用户也能执行docker，需要将用户加入docker组（例如你的登录用户名是togettoyou）

```shell
sudo usermod -aG docker togettoyou # 需要重启生效
```

# 2. 配置 Docker 开机自启

```shell
sudo systemctl enable docker
sudo systemctl start docker
```

查看docker信息

```shell
docker info
```

# 3. 安装 Docker Compose

下载Docker Compose

```shell
curl -L https://get.daocloud.io/docker/compose/releases/download/1.25.4/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
```

配置执行权限

```shell
sudo chmod +x /usr/local/bin/docker-compose
```

检查是否安装成功

```shell
docker-compose -v
```
