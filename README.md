# Wavefront Istio Mixer Adapter

## Usage

### Configuration

1\. Create a copy of the [wavefront/operatorconfig](wavefront/operatorconfig) directory.

2\. If you want the metrics to be published to the Wavefront instance directly, supply
the `direct` params for `wavefront-handler` under `sample_operator_cfg.yaml` like so:

```yaml
...
  connection:
    address: "wavefront-adapter:8080"
  params:
    direct:
      server: https://YOUR-INSTANCE.wavefront.com
      token: YOUR-API-TOKEN
...
```

If you want the metrics to be published to the Wavefront Proxy instead, supply
the `proxy` params like below:

```yaml
...
  connection:
    address: "wavefront-adapter:8080"
  params:
    proxy:
      address: YOUR-PROXY-ADDRESS
...
```

See [config.proto](wavefront/config/config.proto) for the available configuration parameters.

### Deployment

Please follow these steps to configure the Istio Mixer to publish metrics to
Wavefront using this adapter.

1\. Deploy [Istio](https://istio.io/docs/setup/kubernetes/quick-start/).

2\. Deploy the `wavefront-adapter`.

```shell
kubectl apply -f wavefront-adapter.yaml
```

3\. Deploy your copy of the `operatorconfig`.

```shell
kubectl apply -f operatorconfig/
```

You should now be able to see Istio metrics on Wavefront under source `istio`.

## Development

### Setup

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

### Adding Dependencies

To add a dependency, use the following command:

```shell
make vendor-get pkg=<package-uri>
```

For example, you could add a dependency on `istio.io/istio` like so:

```shell
make vendor-get pkg=istio.io/istio
```

### Formatting Code

You could format your code using the following command:

```shell
make format
```

### Building The Docker Image

To build the docker image, use the following command:

```shell
make docker-build
```

### Run The Docker Container

To run the docker container locally, use the following command:

```shell
make docker-run
```

## License
The wavefront-istio-mixer-adapter project is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text.
