# Wavefront by VMware Adapter for Istio

[![Docker Pulls](https://img.shields.io/docker/pulls/vmware/wavefront-adapter-for-istio.svg?logo=docker)](https://hub.docker.com/r/vmware/wavefront-adapter-for-istio/)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE.txt)

<br>
<img alt="Wavefront by VMware" src="docs/images/logo.png">
<br>

Wavefront by VMware Adapter for Istio is an adapter for [Istio](https://istio.io)
to publish metrics to [Wavefront by VMware](https://www.wavefront.com/).

## Quick Start

### Configuration

1\. Create a copy of the [config/operatorconfig/](config/operatorconfig/) directory.

2\. If you want the metrics to be published to the Wavefront instance directly, supply
the `direct` params for `wavefront-handler` under `sample_operator_config.yaml` like so:

```yaml
params:
  direct:
    server: https://YOUR-INSTANCE.wavefront.com
    token: YOUR-API-TOKEN
```

**Note:** Instructions for generating an API token can be found [here](https://docs.wavefront.com/wavefront_api.html#generating-an-api-token).

If you want the metrics to be published to the Wavefront Proxy instead, supply
the `proxy` params like below:

```yaml
params:
  proxy:
    address: YOUR-PROXY-IP:YOUR-PROXY-PORT
```

See the [reference docs](https://preliminary.istio.io/docs/reference/config/policy-and-telemetry/adapters/wavefront/)
for the available configuration parameters.

### Deployment

Please follow these steps to configure the Istio Mixer to publish metrics to
Wavefront using this adapter. These steps must be performed after
deploying [Istio](https://istio.io/docs/setup/kubernetes/quick-start/).

1\. Deploy the `wavefront-adapter.yaml`.

```shell
kubectl apply -f config/wavefront-adapter.yaml
```

2\. Deploy your copy of `operatorconfig`.

```shell
kubectl apply -f your/operatorconfig/
```

You should now be able to see Istio metrics on Wavefront with _istio_ as source.

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md) if you'd like to contribute.

## License

Wavefront by VMware Adapter for Istio is licensed under the Apache License,
Version 2.0. See [LICENSE.txt](LICENSE.txt) for the full license text.
