{{/* Generate rule for tcp metrics */}}
{{- define "rule.tcp" }}
# rule to dispatch tcp metrics to handler wavefront-handler
apiVersion: "config.istio.io/v1alpha2"
kind: rule
metadata:
  name: wavefront-tcp-rule
  namespace: {{ .Values.namespaces.istio }}
spec:
  match: context.protocol == "tcp"
  actions:
  - handler: wavefront-handler.{{ .Values.namespaces.istio }}
    instances:
    - tcpsentbytes.instance.{{ .Values.namespaces.adapter }}
    - tcpreceivedbytes.instance.{{ .Values.namespaces.adapter }}
---
# rule to dispatch tcp connection open metric to handler wavefront-handler
apiVersion: "config.istio.io/v1alpha2"
kind: rule
metadata:
  name: wavefront-tcp-connection-open-rule
  namespace: {{ .Values.namespaces.istio }}
spec:
  match: context.protocol == "tcp" && connection.event == "open"
  actions:
  - handler: wavefront-handler.{{ .Values.namespaces.istio }}
    instances:
    - tcpconnectionsopened.instance.{{ .Values.namespaces.adapter }}
---
# rule to dispatch tcp connection close metric to handler wavefront-handler
apiVersion: "config.istio.io/v1alpha2"
kind: rule
metadata:
  name: wavefront-tcp-connection-close-rule
  namespace: {{ .Values.namespaces.istio }}
spec:
  match: context.protocol == "tcp" && connection.event == "close"
  actions:
  - handler: wavefront-handler.{{ .Values.namespaces.istio }}
    instances:
    - tcpconnectionsclosed.instance.{{ .Values.namespaces.adapter }}
{{- end }}

