# Fetcher Service

This is the Fetcher service

Generated with

```
micro new github.com/sidazhang123/f10-go/srv/fetcher --namespace=sidazhang123.f10 --alias=fetcher --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: sidazhang123.f10.srv.fetcher
- Type: srv
- Alias: fetcher

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
./fetcher-srv
```

Build a docker image
```
make docker
```