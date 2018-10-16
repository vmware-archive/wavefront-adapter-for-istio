# Wavefront by VMware Adapter for Istio

[Helm](https://helm.sh/) is a package manager for Kubernetes. You could use Helm
for installing the Wavefront by VMware adapter on your Kubernetes deployment.

## Quick Start

### Prerequisites

1. Install [Helm](https://docs.helm.sh/using_helm/#installing-helm).

2. Install Tiller via Helm.

```console
$ helm init
```

3. Download and extract [Istio](https://istio.io/docs/setup/kubernetes/download-release/).

4. Install Istio CRDs (Custom Resource Definitions).

```console
$ kubectl apply -f install/kubernetes/helm/istio/templates/crds.yaml
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

It is required that you set the `credentials` for your Wavefront instance. It is
also recommended that you set the `source` attribute to a reasonable value, for
example, to your cluster name.

**Note:** Helm will pick the `direct` credentials by default. If you wish to
ingest metrics via a Proxy, please ensure that the `direct` credentials are
either deleted or commented before deploying.

### Deployment

#### Installation

To install the adapter via Helm, execute the following command.

```console
$ helm install install/wavefront/
```

You should now be able to see Istio metrics on Wavefront under your configured
source (or `istio` by default).

**Note:** On Kubernetes 1.6+, you may encounter the following error if Helm
experiences a problem with RBAC.

```console
$ helm install install/wavefront/
Error: no available release name found
```

To fix the issue, create a Kubernetes Service Account with appropriate
privileges as described in [the Helm documentation](https://docs.helm.sh/using_helm/#tiller-and-role-based-access-control)
and re-install Tiller.

The following example configuration was taken from the [Helm repository](https://github.com/helm/helm/blob/master/docs/rbac.md).

1\. Create a file named `rbac-config.yaml` with the following configuration.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: tiller
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: tiller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - kind: ServiceAccount
    name: tiller
    namespace: kube-system
```

2\. Install the RBAC configuration.

```console
$ kubectl create -f rbac-config.yaml
serviceaccount "tiller" created
clusterrolebinding "tiller" created
```

3\. Reinstall Tiller.

```console
$ helm reset
$ helm init --service-account tiller
```

4\. Install the adapter.

```console
$ helm install install/wavefront/
```

#### Uninstallation

To uninstall the adapter, first identify the Helm release name, like so:

```console
$ helm list
```

Then uninstall it using the following command.

```console
$ helm delete <release-name>
```
