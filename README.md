# Cross AZ Analysis (CAZA)

**WIP**

The goal of this project is to enable the observability of communication within your networks. As engineers we are always trying to improve our infrastructure from a performance and cost standpoint. A large cost for many of us is zone to zone communication in our private networks. In AWS for instance you are charged on what they call "Regional Data Transfer" what this means is that they are billing you based on how the services in your network are communicating, example:

You have a private network with 4 subnets

Subnet A owns the address range 10.1.0.0/16  
Subnet B hwns the address range 10.2.0.0/16  
Subnet C owns the address range 10.4.0.0/16  
Subnet D owns the address range 10.6.0.0/16  

If the services in subnet A are communicating with services in subnet B this is considered a cross zone cost because they belong to different networks. In order to optimize our network communication we need to first see how things are communicating, and currently the options to do that dont really exist.

In order to enable this level of observability I wrote this program "Caza" which utilizes eBPF (extended berkley packet filtering) to read network data from the kernel on TCP_CLOSE events. With this data I can see the source ip and port as well as the destination ip and port which allows us to capture realtime metrics of our network communication.

### Metrics

Metrics are served over port `8080`

We have two types of custom metrics thus far:

- in_network_tx
- out_network_tx

These metrics have an associated "network" label which is a string for the name of the network. In the case of utilizing AWS subnets the "network" label would be the value of the subnets availability zone.

### Usage

**Dependencies**

We do need clang in order to build this program. The Makefile uses `clang-11` but if you are runnning a newer version of clang you can update the Makefile to point to the proper version

`clang-11`  
`libclang-11-dev`  

**Running locally on your host**

```sh
make

sudo ./caza
```

**Building and pushing an image**

```sh
# First ensure you are logged into your docker registry

make image-push
```

**TODO:**

- tests
- improve reported metrics
