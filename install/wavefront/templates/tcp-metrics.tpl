{{/* Generate tcp metrics */}}
{{- define "metrics.tcp" }}
    - name: tcpsentbytes
      instanceName: tcpsentbytes.instance.istio-system
      type: HISTOGRAM
      sample:
        expDecay:
          reservoirSize: 1024
          alpha: 0.015
    - name: tcpreceivedbytes
      instanceName: tcpreceivedbytes.instance.istio-system
      type: HISTOGRAM
      sample:
        expDecay:
          reservoirSize: 1024
          alpha: 0.015
{{- end }}
