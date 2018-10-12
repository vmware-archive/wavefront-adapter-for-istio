# Wavefront by VMware Adapter for Istio

[Helm](https://helm.sh/) is a package manager for Kubernetes. You could use Helm
for installing the Wavefront by VMware adapter on your Kubernetes deployment.

## Quick Start

### Prerequisites

1. Install [Helm](https://docs.helm.sh/using_helm/#installing-helm).

2. Install Tiller via Helm.

```shell
helm init
```

3. Download and extract [Istio](https://istio.io/docs/setup/kubernetes/download-release/).

4. Install Istio CRDs (Custom Resource Definitions).

```shell
kubectl apply -f install/kubernetes/helm/istio/templates/crds.yaml
```

### Configuration

You could configure the adapter installation via `values.yaml`. Below is a list
of configurable parameters.

* **adapter** holds the adapter installation parameters
  * **image** represents the Docker image name
  * **tag** represents the Docker image tag

* **credentials** holds credentials for a Wavefront by VMware instance. One of
  either `direct` or `proxy` parameters must be supplied.
  * **direct** holds credentials for direct ingestion.
    * **server** is the server URL. Ex: https://mydomain.wavefront.com
    * **token** is the API token.
  * **proxy** holds credentials for ingestion via a Proxy.
    * **address** is the proxy address. Ex: 192.168.99.100:2878

* **metrics** holds the metric configuration.
  * **flushInterval** is the metric flush interval.
  * **source** is the source tag for all metrics handled by this adapter.
  * **prefix** is the prefix to prepend all metrics handled by this adapter.
  * **http** is a flag that specifies whether HTTP metrics should be captured.
  * **tcp** is a flag that specifies whether TCP metrics should be captured.

### Deployment

#### Installation

To install the adapter via Helm, execute the following command.

```shell
helm install install/wavefront/
```

#### Uninstallation

To uninstall the adapter, first identify the Helm release name, like so:

```shell
helm list
```

Then uninstall it using the following command.

```shell
helm delete <release-name>
```
