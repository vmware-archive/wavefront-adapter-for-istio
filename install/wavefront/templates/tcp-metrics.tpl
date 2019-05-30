{{/* Generate tcp metrics */}}
{{- define "metrics.tcp" }}
    - name: tcpsentbytes
      instanceName: tcpsentbytes.instance.{{ .Values.adapterNamespace }}
      type: HISTOGRAM
      sample:
        expDecay:
          reservoirSize: 1024
          alpha: 0.015
    - name: tcpreceivedbytes
      instanceName: tcpreceivedbytes.instance.{{ .Values.adapterNamespace }}
      type: HISTOGRAM
      sample:
        expDecay:
          reservoirSize: 1024
          alpha: 0.015
{{- end }}
