
<p align="center">

[![build status](https://github.com/werbenhu/registry/workflows/Go/badge.svg)](https://github.com/werbenhu/registry/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/werbenhu/registry.svg)](https://pkg.go.dev/github.com/werbenhu/registry)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/werbenhu/registry/issues)
[![Mit License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://pkg.go.dev/github.com/werbenhu/registry)
</p>

[English](README.md) | [简体中文](README-CN.md)

# Registry
**一个轻量的服务注册中心，使用一致性哈希算法进行服务发现。**

## 什么是一致性哈希

一致性哈希(Consistent Hash)是为了解决由于分布式系统中节点的增加或减少而带来的大量失效问题，它可以有效地降低这种失效影响，从而提高分布式系统的性能和可用性。

### 普通哈希的问题

普通哈希函数是 key % n，其中 n 是服务器数量。它有两个主要缺点：
1. 不能水平扩展，或者换句话说，不具备分区容错性。当添加新服务器时，所有现有的映射都会被破坏。这可能会引入痛苦的维护工作和系统停机时间。
2. 可能不能实现负载均衡。如果数据不是均匀分布的，这可能会导致一些服务器过热饱和，而其他服务器则处于空闲状态并几乎为空。

问题2可以通过先对键进行哈希，然后哈希(key) % n，以便哈希键更有可能被均匀分布来解决。但是，这不能解决问题1。我们需要找到一个可以分配key并且不依赖于n的解决方案。

### 一致性哈希

关于环，添加删除节点，还有虚拟节点等，参考：[一致性哈希的简单认识](https://baijiahao.baidu.com/s?id=1735480432495470467&wfr=spider&for=pc)

## 如何运行

### 编译
```sh
cd cmd 
go build -o registry
```

### 使用方法
```
  -id string
        服务ID，不能为空
  -bind string
        用于注册服务的地址 (默认为":7370")。
  -bind-advertise string
        服务将向其他服务公布的地址 (默认为":7370")。
  -addr string
        用于服务发现的地址 (默认为":9800")。
  -advertise string
        服务向客户端公布的地址以供服务发现 (默认为":9800")。
  -registries string
        注册中心服务器地址，可以为空，多个地址用逗号分隔。
  
```
## 启动注册中心服务器

要启动注册中心服务器，请按照以下步骤进行操作：

1. 根据实际情况确定所需节点数。
2. 执行以下命令启动节点：

``` sh
# 启动第一个节点
./registry -bind=":7370" \
     -bind-advertise="172.16.3.3:7370" \
     -id=service-1 \
     -addr=":9800" \
     -advertise="172.16.3.3:9800"

# 启动第二个节点
# 第二个节点有一个额外的参数-registries="172.16.3.3:7370"，
# 因为第二个节点需要向第一个节点注册
./registry -bind=":7371" \
     -bind-advertise="172.16.3.3:7371" \
     -id=service-2 \
     -registries="172.16.3.3:7370" \
     -addr=":9801" \
     -advertise="172.16.3.3:9801"
```

注意：如果存在防火墙，请确保同时打开TCP和UDP端口。


## 注册服务

使用以下代码片段注册服务：

```
// Create a new registration object
r := register.New(id, bind, advertise, registries, group, addr)

// Start the registration
err = r.Start()
if err != nil {
	panic(err)
}
```

Parameters:

参数：
- id：服务ID。
- bind：用于将服务注册到注册表服务器的地址。
- advertise：服务将向注册表服务器公布的地址。可用于基本的NAT遍历，其中内部IP：端口和外部IP：端口均已知。
- registries：注册表服务器的地址。如果有多个，请使用逗号分隔，例如"192.168.1.101:7370,192.168.1.102:7370"。
- group：当前服务所属的组名称。
- addr：当前服务向客户端提供的地址。例如，如果当前服务是一个HTTP服务器，则地址是172.16.3.3:80，即HTTP的地址。


## 服务发现
### 用法
```
// 您可以选择已启动的任何一个注册服务器。
registryAddr := "172.16.3.3:9801"
group := "test-group"

// 创建新的RpcClient
client, err := client.NewRpcClient(registryAddr)
if err != nil {
	panic(err)
}

// 使用一致性哈希根据用户ID分配服务
service, err := client.Match(groupName, "user-id-1")
if err != nil {
	panic(err)
}

log.Printf("[INFO] Matched key: %s, Service ID: %s, Service Address: %s\n", key, service.Id, service.Addr)

// 获取该组所有服务
allService, err := client.Members(group)
if err != nil {
      log.Printf("[ERROR] Failed to get all services: %s\n", err)
}
log.Printf("[INFO] All services: %+v\n", allService)
```

## 示例

### 注册两个 Web 服务
```sh
# 注册第一个 Web 服务
cd examples/service
go build -o webservice webservice.go 
./webservice \
	-group=webservice-group \
	-id=webserver1 \
	-registries=172.16.3.3:7370 \
	-bind=":8370"
	-addr="172.16.3.3:8080"

# 注册第二个 Web 服务
cd examples/service
./webservice \
	-group=webservice-group \
	-id=webserver2 \
	-registries=172.16.3.3:7370 \
	-bind=":8371" \
	-addr="172.16.3.3:8081"
```

### 客户端发现服务
```sh
cd examples/client
go build -o client main.go
./client
```

## 贡献
欢迎贡献和反馈！提出问题或提出功能请求请提 [issue](https://github.com/werbenhu/registry/issues) .


