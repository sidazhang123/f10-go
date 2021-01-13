# Feed Service

This is the Feed service

Generated with

```
micro new f10-go/srv/feed --namespace=sidazhang123.f10 --alias=feed --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: sidazhang123.f10.srv.feed
- Type: srv
- Alias: feed

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend etcd.

```
# install etcd
brew install etcd

# run etcd
etcd
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./feed-srv
```

Build a docker image
```
make docker
```