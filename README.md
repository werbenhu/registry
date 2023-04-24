
<p align="center">

[![build status](https://github.com/werbenhu/registry/workflows/Go/badge.svg)](https://github.com/werbenhu/registry/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/werbenhu/registry.svg)](https://pkg.go.dev/github.com/werbenhu/registry)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/werbenhu/registry/issues)
[![Mit License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://pkg.go.dev/github.com/werbenhu/registry)
</p>

[English](README.md) | [简体中文](README-CN.md)

# Registry
**Registry is a simple service registry that uses the consistent hashing algorithm for service discovery.**

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

To start a registry server, follow these steps:

1. Determine the number of nodes required based on your actual situation.
2. Execute the following commands to start the nodes:

``` sh
# Starting the first node
./registry -bind=":7370" \
     -bind-advertise="172.16.3.3:7370" \
     -id=service-1 \
     -addr=":9800" \
     -advertise="172.16.3.3:9800"

# Starting the second node
# The second one has an additional parameter -registries="172.16.3.3:7370",
# because the second node needs to register with the first one
./registry -bind=":7371" \
     -bind-advertise="172.16.3.3:7371" \
     -id=service-2 \
     -registries="172.16.3.3:7370" \
     -addr=":9801" \
     -advertise="172.16.3.3:9801"
```

Note: If there is a firewall, make sure to open both TCP and UDP on advertise ports.


## Register services

Use the following code snippet to register services:

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

- `id`: Service ID.
- `bind`: Address used to register the service to the registry server.
- `advertise`: Address that the service will advertise to the registry server. Can be used for basic NAT traversal where both the internal IP:port and external IP:port are known.
- `registries`: Addresses of the registry server(s). If there are more than one, separate them with commas, such as "192.168.1.101:7370,192.168.1.102:7370".
- `group`: Group name the current service belongs to.
- `addr`: Address currently provided by this service to the client. For example, if the current service is an HTTP server, the address is 172.16.3.3:80, which is the address that HTTP listens to.


## Service Discovery
### Usage
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
if err != nil {
	panic(err)
}

log.Printf("[INFO] Matched key: %s, Service ID: %s, Service Address: %s\n", key, service.Id, service.Addr)

// Get all services of the group
allService, err := client.Members(group)
if err != nil {
      log.Printf("[ERROR] Failed to get all services: %s\n", err)
}
log.Printf("[INFO] All services: %+v\n", allService)
```

## Examples

### Register two web services.
```sh
# Register the first web service.
cd examples/service
go build -o webservice webservice.go 
./webservice \
	-group=webservice-group \
	-id=webserver1 \
	-registries=172.16.3.3:7370 \
	-bind=":8370" \
	-advertise="172.16.3.3:8370" \
	-addr="172.16.3.3:8080"

# Register the second web service.
cd examples/service
./webservice \
	-group=webservice-group \
	-id=webserver2 \
	-registries=172.16.3.3:7370 \
	-bind=":8371" \
	-advertise="172.16.3.3:8371" \
	-addr="172.16.3.3:8081"
```

### Client discovery service
```sh
cd examples/client
go build -o client main.go
./client
```

## Contributions
Contributions and feedback are both welcomed and encouraged! Open an [issue](https://github.com/werbenhu/registry/issues) to report a bug, ask a question, or make a feature request.

