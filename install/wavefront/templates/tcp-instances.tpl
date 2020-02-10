{{/* Generate instances for tcp metrics */}}
{{- define "instances.tcp" }}
# tcpsentbytes instance for template metric
apiVersion: "config.istio.io/v1alpha2"
kind: instance
metadata:
  name: tcpsentbytes
  namespace: {{ .Values.namespaces.adapter }}
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
  namespace: {{ .Values.namespaces.adapter }}
spec:
  template: metric
  params:
    value: connection.received.bytes | 0
    dimensions:
      {{- template "attributes.service" }}
    monitored_resource_type: '"UNSPECIFIED"'
---
# tcpconnectionsopened instance for template metric
apiVersion: "config.istio.io/v1alpha2"
kind: instance
metadata:
  name: tcpconnectionsopened
  namespace: {{ .Values.namespaces.adapter }}
spec:
  template: metric
  params:
    value: 1
    dimensions:
      {{- template "attributes.service" }}
    monitored_resource_type: '"UNSPECIFIED"'
---
# tcpconnectionsclosed instance for template metric
apiVersion: "config.istio.io/v1alpha2"
kind: instance
metadata:
  name: tcpconnectionsclosed
  namespace: {{ .Values.namespaces.adapter }}
spec:
  template: metric
  params:
    value: 1
    dimensions:
      {{- template "attributes.service" }}
    monitored_resource_type: '"UNSPECIFIED"'
{{- end }}
