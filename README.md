# Wavefront Istio Mixer Adapter

## Usage

### Configuration

1\. Create a copy of the [config/operatorconfig/](config/operatorconfig/) directory.

2\. If you want the metrics to be published to the Wavefront instance directly, supply
the `direct` params for `wavefront-handler` under `sample_operator_config.yaml` like so:

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
kubectl apply -f config/wavefront-adapter.yaml
```

3\. Deploy your copy of the `operatorconfig`.

```shell
kubectl apply -f your/operatorconfig/
```

You should now be able to see Istio metrics on Wavefront under source `istio`.

## Contributing

Please check out [CONTRIBUTING.md](CONTRIBUTING.md) if you'd like to contribute.

## License

Wavefront Istio Mixer Adapter is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text.
