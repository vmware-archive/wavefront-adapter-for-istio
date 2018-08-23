# Wavefront Istio Mixer Adapter

## Development

### Setup

1\. Install [Golang](https://golang.org/dl/).

2\. Set `GOPATH` and `GOBIN` like so:

```shell
export GOPATH=~/go
export GOBIN=$GOPATH/bin
```

3\. Install the development tools like so:

```shell
make setup
```

Run `make help` to get a list of all available targets.

### Adding Dependencies

To add a dependency, use the following command:

```shell
make vendor-get <package-uri>
```

For example, you could add a dependency on `istio.io/istio` like so:

```shell
make vendor-get istio.io/istio
```

### Formatting Code

You could format your code using the following command:

```shell
make format
```
