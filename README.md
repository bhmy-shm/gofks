
# Gofks 微服务脚手架

Gofks 是一个由Go编写的微服务脚手架。集成多个开源框架优秀思想。方便一键式上手编写业务代码，快速开发。

**实现功能:**
1. 命令行工具支持(一键生成项目文件 webApi | webSocket | rpcServer | plugins )
2. 兼容 gin 框架，极简的api调用
3. 兼容 gRpc，支持中间件拦截器，方便扩展。内建限流、超时、链路追踪、自动恢复
4. 兼容 plugins 阻塞式的插件类服务 （类似cron）
5. 兼容 websocket 实现高并发可控的长连接控制，及路由管理调度gRpc
6. 自定义可靠的、易扩展的 消息传输协议
7. otel 链路追踪（ rpc | websocket | db | redis ）
8. 配置热更新


## Getting started


With [Go module](https://github.com/golang/go/wiki/Modules) support, simply add the following import

```
import "github.com/bhmy-shm/gofks"
```

run the following Go command to install the `gofks` package:

```sh
$ go get -u github.com/bhmy-shm/gofks
```
## Quick started
基于go-zero 开源框架二开实现
1. 安装gofkctl 工具
```sh
$ go install github.com/bhmy-shm/gofk
```

2. 快速生成 rpc 服务
```sh
$ gofkctl rpc new hello-rpc
```

3. 快速生成 web-api 服务
```sh
$ gofkctl api new hello-web --rpc=hello-rpc
```

4. 快速生成 web-api + websocket 服务
```sh
$ gofkctl api new hello-web --rpc=hello-rpc --websocket=true
```

5. 快速生成 plugins 服务
```sh
$ gofkctl plugins new hello-plugins --server=first,secod,third -registry=ture  
```

## example Directory

1. [web-api](./example/api/README.md)
2. [webSocket](./example/websocket/controls/handler_ws.go)
3. [rpc-server](./example/rpc/account-rpc/main.go)
3. [plugins](./example/plugins/manager.go)
4. [model-db](./example/model/account/account.go)

## 项目中 rpc 与 api-websocket 关联创建

1. 创建项目目录，编写一个proto文件
```protobuf
syntax = "proto3";

package user;
option go_package = "./user";

message Status {
  uint64 code = 1 [json_name = "code,required"];
  string reason = 2;
  string message = 3;
  map<string, string> metadata = 4;
};

message PageParam {
  int32 pageNum = 1; //当前页
  int32 pageSize = 2; //查询数量
}

service UserClient {
  //用户配置
  rpc Login(LoginReq)  returns (LoginResp);
  rpc UpPass(UpPassReq) returns (UpPassResp);
}
```

2. 基于proto文件生成 rpc项目
```shell
$ cd mall/user/rpc
$ gofkctl rpc protoc user.proto \
  --go_out=./types --go-grpc_out=./types \
  --name=user-rpc \
  --port=9091 
```

3. 基于rpc项目 生成 websocket 和 api服务
```shell
$ gofkctl api new v1 --websocket --rpc=user-rpc
```



