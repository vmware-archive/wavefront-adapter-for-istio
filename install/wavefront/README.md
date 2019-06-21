# Wavefront by VMware Adapter for Istio

[Helm](https://helm.sh/) is a package manager for Kubernetes. You could use Helm
for installing the Wavefront by VMware adapter on your Kubernetes deployment.

## Quick Start

### Prerequisites

To deploy this adapter, you will need a cluster with the following minimum setup.

* Kubernetes v1.10.0
* Istio v1.0.3
* Helm v2.10.0

### Helm Setup

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

1\. Download the Helm chart configuration and extract it.

```console
$ curl -LO https://github.com/vmware/wavefront-adapter-for-istio/releases/download/0.1.1/wavefront-0.1.1.tgz
$ tar -zxvf wavefront-0.1.1.tgz
```

2\. The configuration used per Helm deployment is specified in the `values.yaml`
file.

**Note:** Helm will pick the `direct` credentials by default. If you wish to
ingest metrics via a Proxy, please ensure that the `direct` credentials are
either deleted or commented before deploying.

If you want the metrics to be published to the Wavefront instance directly,
supply the `direct` params like so:

```yaml
credentials:
  direct:
    server: https://YOUR-INSTANCE.wavefront.com
    token: YOUR-API-TOKEN
```

Instructions for generating an API token can be found in the Wavefront by VMware
[docs](https://docs.wavefront.com/wavefront_api.html#generating-an-api-token).

If you want the metrics to be published to the Wavefront Proxy instead, supply
the `proxy` params like below:

```yaml
credentials:
  proxy:
    address: YOUR-PROXY-IP:YOUR-PROXY-PORT
```

3\. It is recommended that you update the `source` attribute to a reasonable
value, for example, to your cluster name.

```yaml
metrics:
  source: my-cluster
```

#### Configuration Parameters

See below for the list of available configuration parameters.

| Parent      | Parameter     | Description                                        |
| ----------- | ------------- | -------------------------------------------------- |
| adapter     | image         | The Docker image name                              |
|             | tag           | The Docker image tag                               |
| credentials | direct        | Credentials for direct ingestion                   |
|             | proxy         | Credentials for ingestion via a Proxy              | 
| direct      | server        | The Server URL. Ex: https://mydomain.wavefront.com |
|             | token         | The API token                                      |
| proxy       | address       | The Proxy address. Ex: 192.168.99.100:2878         |
| metrics     | flushInterval | The metric flush interval                          |
|             | source        | The source tag for all metrics                     |
|             | prefix        | The prefix to prepend all metrics with             |
|             | http          | Specify whether HTTP metrics should be captured    |
|             | tcp           | Specify whether TCP metrics should be captured     |
| namespaces  | adapter       | The namespace to create adapter objects in         |
|             | istio         | The namespace Istio has been installed to          |

### Deployment

#### Installation

To install the adapter via Helm, execute the following command.

```console
$ helm install wavefront/
```

You should now be able to see Istio metrics on Wavefront under your configured
source (or `istio` by default).

##### RBAC Issues On Helm

On Kubernetes 1.6+, you may encounter the following error if Helm experiences a
problem with RBAC.

```console
$ helm install wavefront/
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
$ helm install wavefront/
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
