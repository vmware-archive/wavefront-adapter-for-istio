{{/* Generate instances for tcp metrics */}}
{{- define "instances.tcp" }}
# tcpsentbytes instance for template metric
apiVersion: "config.istio.io/v1alpha2"
kind: instance
metadata:
  name: tcpsentbytes
  namespace: {{ .Values.adapterNamespace }}
spec:
  template: metric
  params:
    value: connection.sent.bytes | 0
    dimensions:
      {{- template "attributes.service" }}
    monitored_resource_type: '"UNSPECIFIED"'
---
# tcpreceivedbytes instance for template metric
apiVersion: "config.istio.io/v1alpha2"
kind: instance
metadata:
  name: tcpreceivedbytes
  namespace: {{ .Values.adapterNamespace }}
spec:
  template: metric
  params:
    value: connection.received.bytes | 0
    dimensions:
      {{- template "attributes.service" }}
    monitored_resource_type: '"UNSPECIFIED"'
{{- end }}
