# Accumulator Service

This is the Accumulator service

Generated with

```
micro new --namespace=sidazhang123.f10.srv --type=service accumulator
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: sidazhang123.f10.srv.service.accumulator
- Type: service
- Alias: accumulator

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
./accumulator-service
```

Build a docker image
```
make docker
```