# exposing_intranet

内网穿透服务，server 模块同时提供服务端和客户端运行能力。

## 目录结构

```text
server/
├── app/                 服务端业务分层
├── client/              客户端 API、运行时和传输层
├── cmd/                 Cobra 命令入口
├── module/              服务端技术模块
├── tools/               公共协议和工具模块
└── website/             配置、数据库和静态资源
```

## 运行

服务端默认命令：

```shell
go run .
```

客户端使用服务端程序的 `client` 子命令：

```shell
go run . client
go run . client auth <username:password>
go run . client map add <server_port> <local_addr> [comment]
go run . client map del <server_port>
go run . client map list
```

## 配置

- 服务端配置：`website/configs/config.yaml`
- 客户端配置：`website/configs/config.client.yaml`
- 客户端 symbol：`website/configs/client.key`

客户端配置可通过 `EXPOSING_INTRANET_CLIENT_CONFIG` 覆盖，symbol 路径可通过
`EXPOSING_INTRANET_SYMBOL_PATH` 覆盖。

## 构建

```shell
GOOS=linux CGO_ENABLED=0 GOARCH=amd64 \
  go build -ldflags="-s -w" -o ./build/linux_service main.go
```

构建后的同一个程序同时支持：

```shell
./build/linux_service
./build/linux_service client
```

同步发布目录：

```shell
./sync_build.sh all
```
