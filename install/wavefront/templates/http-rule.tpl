{{/* Generate rule for http metrics */}}
{{- define "rule.http" }}
# rule to dispatch to handler wavefront-handler
apiVersion: "config.istio.io/v1alpha2"
kind: rule
metadata:
  name: wavefront-http-rule
  namespace: {{ .Values.namespaces.istio }}
spec:
  actions:
  - handler: wavefront-handler.{{ .Values.namespaces.istio }}
    instances:
    - requestsize.instance.{{ .Values.namespaces.adapter }}
    - requestcount.instance.{{ .Values.namespaces.adapter }}
    - requestduration.instance.{{ .Values.namespaces.adapter }}
    - responsesize.instance.{{ .Values.namespaces.adapter }}
{{- end }}
