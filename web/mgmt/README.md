# Regex Service

This is the Regex service

Generated with

```
micro new github.com/sidazhang123/f10-go/web/regex --namespace=sidazhang123.f10 --alias=regex --type=web
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: sidazhang123.f10.web.regex
- Type: web
- Alias: regex

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
./regex-web
```

Build a docker image
```
make docker
```