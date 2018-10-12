# Wavefront by VMware Adapter for Istio

[![Docker Pulls](https://img.shields.io/docker/pulls/vmware/wavefront-adapter-for-istio.svg?logo=docker)](https://hub.docker.com/r/vmware/wavefront-adapter-for-istio/)
[![License](https://img.shields.io/badge/license-Apache--2.0-green.svg)](LICENSE)

<img alt="Wavefront by VMware" src="docs/images/logo.png">

Wavefront by VMware Adapter for Istio is an adapter for [Istio](https://istio.io)
to publish metrics to [Wavefront by VMware](https://www.wavefront.com/).

**Note:** This adapter is currently experimental. Therefore, caution should be
taken before using it in production environments.

## Quick Start

### Configuration

1\. Download the configuration from the [releases page](https://github.com/vmware/wavefront-adapter-for-istio/releases)
and extract it.

```shell
curl -LO https://github.com/vmware/wavefront-adapter-for-istio/releases/download/0.1.0/config.tar.gz
tar -zxvf config.tar.gz
```

2\. If you want the metrics to be published to the Wavefront instance directly,
supply the `direct` params for `wavefront-handler` under
`sample_operator_config.yaml` like so:

```yaml
params:
  direct:
    server: https://YOUR-INSTANCE.wavefront.com
    token: YOUR-API-TOKEN
```

Instructions for generating an API token can be found in the Wavefront by VMware
[docs](https://docs.wavefront.com/wavefront_api.html#generating-an-api-token).

If you want the metrics to be published to the Wavefront Proxy instead, supply
the `proxy` params like below:

```yaml
params:
  proxy:
    address: YOUR-PROXY-IP:YOUR-PROXY-PORT
```

3\. It is recommended that you update the `source` attribute to a reasonable
value, for example, to your cluster name.

```yaml
params:
  ...
  source: my-cluster
```

See the [reference docs](https://istio.io/docs/reference/config/policy-and-telemetry/adapters/wavefront/)
for the available configuration parameters.

### Deployment

Execute the following command to configure the Istio Mixer to publish metrics to
Wavefront using this adapter. This step must be performed after deploying
[Istio](https://istio.io/docs/setup/kubernetes/quick-start/).

```shell
kubectl apply -f config/
```

You should now be able to see Istio metrics on Wavefront under your configured
source (or `istio` by default).

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md) if you'd like to contribute.

## License

Wavefront by VMware Adapter for Istio is licensed under the Apache License,
Version 2.0. See [LICENSE](LICENSE) for the full license text. Also, see the
[open_source_licenses](open_source_licenses) file for the full license text from
the packages used in this project.
