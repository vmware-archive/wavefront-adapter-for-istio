# Build

## Setup

1\. Install [Golang](https://golang.org/dl/).

2\. Install [Docker](https://github.com/istio/istio/wiki/Dev-Guide#setting-up-docker).

3\. Set `GOPATH` and `GOBIN` like so:

```shell
export GOPATH=~/go
export GOBIN=$GOPATH/bin
```

4\. Install the development tools like so:

```shell
make setup
```

Run `make help` to get a list of all available targets.

## Adding Dependencies

To add a dependency, use the following command:

```shell
make vendor-get pkg=<package-uri>
```

For example, you could add a dependency on `istio.io/istio` like so:

```shell
make vendor-get pkg=istio.io/istio
```

## Formatting Code

You could format your code using the following command:

```shell
make format
```

## Building The Docker Image

To build the docker image, use the following command:

```shell
make docker-build
```

## Running The Docker Container

To run the docker container locally, use the following command:

```shell
make docker-run
```

## Test

To run the unit tests, use the following command:

```shell
make test
```

## Dry-Running Helm

To dry-run the helm chart, use the following command:

```shell
make helm-print
```

## Generating Helm Manifest

To generate the `install/config.yaml` manifest for the Helm chart, use the
following command:

```shell
make helm-generate
```

## Packing Helm Configuration

To pack Helm configuration files for releases, use the following command:

```shell
make helm-pack
```
