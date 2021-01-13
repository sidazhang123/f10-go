# Processor Service

This is the Processor service

Generated with

```
micro new github.com/sidazhang123/f10-go/srv/processor --namespace=sidazhang123.f10 --alias=processor --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: sidazhang123.f10.srv.processor
- Type: srv
- Alias: processor

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend consul.

```
# install consul
brew install consul

# run consul
consul agent -dev
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./processor-srv
```

Build a docker image
```
make docker
```