# srouter
一个用纯Go语言写的一致性哈希路由服务器，支持自动发现。

#### 什么是一致性哈希

> 百度百科：一种特殊的哈希算法，目的是解决分布式缓存的问题。在移除或者添加一个服务器时，能够尽可能小地改变已存在的服务请求与处理请求服务器之间的映射关系。一致性哈希解决了简单哈希算法在分布式哈希表中存在的动态伸缩等问题。

假设有1000万个用户，100个服务器node，请设计一种算法合理地将用户分配到这些服务器上。普通的哈希算法是将1000万个用户各自的userid计算出hash值，取余100，然后选择对应编号的服务器。由于该算法使用节点数取余的方法，强依赖node的数目，当node数发生变化的时候，比如100个服务器有5个宕机了，现在只剩95个。这时候每个用户都需要重新分配服务器，1000万个用户的id计算出hash值，取余95，大部分用户所属的服务都要变更。

> 一致性哈希主要解决的问题就是：当node数发生变化时，能够尽量少的移动数据。

#### 如何运行

##### 编译路由服务
```sh
cd cmd 
go build -o srouter.exe
```

##### Usage命令
```
  -addr string
        服务发现通信的地址 (default ":7370")
  -advertise string
        对外公布的服务发现通信的地址 (default ":7370")
  -api-addr string
        查询服务器的地址 (default ":8080")
  -id string
        服务ID，不能为空
  -routers string
        路由服务器地址，如果是第一个可以为空，多个用逗号隔开
  -service string
        对外公布的查询服务器的地址
```

##### 启动路由服务器
``` sh
# 这里演示启动2个，启动数量可以自己根据实际情况定
# 启动第1个
./srouter.exe -addr=":7370" `
     -advertise="172.16.3.3:7370" `
     -id=router-service-1 `
     -api-addr=":9000" `
     -service="172.16.3.3:9000"

# 启动第2个
# 第2个多一个参数-routers="172.16.3.3:7370"
# 这里是需要将第2个注册到第1个去
./srouter.exe -addr=":7371" `
     -advertise="172.16.3.3:7371" `
     -id=router-service-2 `
     -routers="172.16.3.3:7370" `
     -api-addr=":9001" `
     -service="172.16.3.3:9001"
```
##### 注册2个web服务
```sh
# 注册第1个web服务, 
# 注册的服务组是webservice-group
# 注册的服务ID是webserver1
# 注册的服务地址是172.16.3.3:8000
cd examples/service1
go build -o webservice1.exe webservice1.go 
./webservice1.exe

# 注册第2个web服务
# 注册的服务组是webservice-group
# 注册的服务ID是webserver2
# 注册的服务地址是172.16.3.3:8001
cd examples/service2
go build -o webservice2.exe webservice2.go
./webservice2.exe
```

##### 客户端使用一致性哈希选择服务
```sh
cd examples/client
go build -o client.exe main.go
./client.exe
```

