{{/* Generate rule for http metrics */}}
{{- define "rule.http" }}
# rule to dispatch to handler wavefront-handler
apiVersion: "config.istio.io/v1alpha2"
kind: rule
metadata:
  name: wavefront-http-rule
  namespace: {{ .Values.istioNamespace }}
spec:
  actions:
  - handler: wavefront-handler.{{ .Values.istioNamespace }}
    instances:
    - requestsize.instance.{{ .Values.adapterNamespace }}
    - requestcount.instance.{{ .Values.adapterNamespace }}
    - requestduration.instance.{{ .Values.adapterNamespace }}
    - responsesize.instance.{{ .Values.adapterNamespace }}
{{- end }}
