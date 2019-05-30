{{/* Generate instances for http metrics */}}
{{- define "instances.http" }}
# requestsize instance for template metric
apiVersion: "config.istio.io/v1alpha2"
kind: instance
metadata:
  name: requestsize
  namespace: {{ .Values.adapterNamespace }}
spec:
  template: metric
  params:
    value: request.total_size | 0
    dimensions:
      {{- template "attributes.service" }}
      response_code: response.code | 200
    monitored_resource_type: '"UNSPECIFIED"'
---
# requestcount instance for template metric
apiVersion: "config.istio.io/v1alpha2"
kind: instance
metadata:
  name: requestcount
  namespace: {{ .Values.adapterNamespace }}
spec:
  template: metric
  params:
    value: 1
    dimensions:
      {{- template "attributes.service" }}
      response_code: response.code | 200
    monitored_resource_type: '"UNSPECIFIED"'
---
# requestduration instance for template metric
apiVersion: "config.istio.io/v1alpha2"
kind: instance
metadata:
  name: requestduration
  namespace: {{ .Values.adapterNamespace }}
spec:
  template: metric
  params:
    value: response.duration | "0ms"
    dimensions:
      {{- template "attributes.service" }}
      response_code: response.code | 200
    monitored_resource_type: '"UNSPECIFIED"'
---
# responsesize instance for template metric
apiVersion: "config.istio.io/v1alpha2"
kind: instance
metadata:
  name: responsesize
  namespace: {{ .Values.adapterNamespace }}
spec:
  template: metric
  params:
    value: response.total_size | 0
    dimensions:
      {{- template "attributes.service" }}
      response_code: response.code | 200
    monitored_resource_type: '"UNSPECIFIED"'
{{- end }}
