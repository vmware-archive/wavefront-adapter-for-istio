{{/* Generate rule for tcp metrics */}}
{{- define "rule.tcp" }}
# rule to dispatch tcp metrics to handler wavefront-handler
apiVersion: "config.istio.io/v1alpha2"
kind: rule
metadata:
  name: wavefront-tcp-rule
  namespace: istio-system
spec:
  match: context.protocol == "tcp"
  actions:
  - handler: wavefront-handler.istio-system
    instances:
    - tcpsentbytes
    - tcpreceivedbytes
{{- end }}
