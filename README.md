# Terraform Provider for SAP DI

This is the repository for the Terraform SAP DI provider, which one can use with Terraform to work with SAP DI.

**Disclaimer**: This provider is currently in alpha stage and is not yet ready for production use.
This is a community project and not officially supported by SAP.
Proceed with caution.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To compile the provider fo another OS/Architecture, run `GOOS=<OS> GOARCH=<ARCH> go build .`.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.
Currently, these tests only test against a local mock server which can be started manually with `make start-mock-server`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
