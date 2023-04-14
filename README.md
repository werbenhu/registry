# Registry
**A simple registry server to discover your services, it uses consistent hashing algorithm for service discovery.**

## What is consistent hashing

> Consistent hashing is a hashing technique that performs really well when operated in a dynamic environment where the distributed system scales up and scales down frequently. 

### The problem of naive hashing function

A naive hashing function is key % n where n is the number of servers.
It has two major drawbacks:
1. NOT horizontally scalable, or in other words, NOT partition tolerant. When you add new servers, all existing mapping are broken. It could introduce painful maintenance work and downtime to the system.
2. May NOT be load balanced. If the data is not uniformly distributed, this might cause some servers to be hot and saturated while others idle and almost empty.

Problem 2 can be resolved by hashing the key first, hash(key) % n, so that the hashed keys will be likely to be distributed more evenly. But this can't solve the problem 1. We need to find a solution that can distribute the keys and is not dependent on n.

### Consistent Hashing
Consistent Hashing allows distributing data in such a way that minimize reorganization when nodes are added or removed, hence making the system easier to scale up or down.

The key idea is that it's a distribution scheme that DOES NOT depend directly on the number of servers.

In Consistent Hashing, when the hash table is resized, in general only k / n keys need to be remapped, where k is the total number of keys and n is the total number of servers.

When a new node is added, it takes shares from a few hosts without touching other's shares
When a node is removed, its shares are shared by other hosts.

## Getting started

### Build registry
```sh
cd cmd 
go build -o registry
```

### Usage
```
   -id string
        Service ID, cannot be empty
  -bind string
        The address used to register the service (default ":7370").
  -bind-advertise string
        The address will advertise to other services (default ":7370").
  -addr string
        The address used for service discovery (default ":9800").
  -advertise string
        The address will advertise to client for service discover (default ":9800").
  -registries string
        Registry server addresses, it can be empty, and multiples are separated by commas.
  
```
## Starting registry server

Here is an example, start 2 registry server nodes, the number of starts can be determined according to the actual situation.


``` sh
# starting the first node
./registry -bind=":7370" \
     -bind-advertise="172.16.3.3:7370" \
     -id=service-1 \
     -addr=":9800" \
     -advertise="172.16.3.3:9800"

# starting the second node
# The second one has one more parameter -registries="172.16.3.3:7370",
# this is because the second one needs to be registered to the first one
./registry -bind=":7371" \
     -bind-advertise="172.16.3.3:7371" \
     -id=service-2 \
     -registries="172.16.3.3:7370" \
     -addr=":9801" \
     -advertise="172.16.3.3:9801"
```

## Register services

```
// id: 
//    service ID
// bind: 
//    The address used to register the service to registry server.
//    If there is a firewall, please remember that the port needs to open both tcp and udp.
// advertise: 
//    The address that the service will advertise to registry server. 
//    Can be used for basic NAT traversal where both the internal ip:port and external ip:port are known.
// registries:
//    The addresses of the registry server, if there are more than one, separate them with commas, 
//    such as "192.168.1.101:7370,192.168.1.102:7370"
// group: 
//    Group name the current service belongs to.
// addr: 
//    The address currently provided by this service to the client, 
//    for example, the current service is an http server, 
//    that is the address 172.16.3.3:80 that http listens to.

r := register.New(id, bind, advertise, registries, group, addr)
err = r.Start()
if err != nil {
	panic(err)
}
```


## Service Discovery
```
// You can choose any one of the registered servers.
registryAddr := "172.16.3.3:9801"
group := "test-group"

// Create a new RpcClient
client, err := client.NewRpcClient(registryAddr)
if err != nil {
	panic(err)
}

// Use consistent hash to assign services based on user ID
service, err := client.Match(groupName, "user-id-1")
log.Printf("[INFO] match key:%s, serviceId:%s, serviceAddr:%s\n", key, service.Id, service.Addr)

// Get all services of the group
allService, err := client.Members(group)
if err != nil {
	log.Printf("[ERROR] get all service err:%s\n", err)
}
log.Printf("[INFO] all service:%+v\n", allService)
```

## Examples

### Register two web services.
```sh
# Register the first web service, 
# The service group is webservice-group,
# The service ID is webserver1
# The web service address is 172.16.3.3:8080
cd examples/service1
go build -o webservice1 webservice1.go 
./webservice1

# Register the second web service, 
# The service group is webservice-group,
# The service ID is webserver2
# The web service address is 172.16.3.3:8081
cd examples/service2
go build -o webservice2 webservice2.go
./webservice2
```

### Client discovery service
```sh
cd examples/client
go build -o client main.go
./client
```

## Contributions
Contributions and feedback are both welcomed and encouraged! Open an [issue](https://github.com/werbenhu/registry/issues) to report a bug, ask a question, or make a feature request.

