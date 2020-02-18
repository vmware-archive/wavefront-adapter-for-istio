{{/* Generate tcp metrics */}}
{{- define "metrics.tcp" }}
    - name: tcpsentbytes
      instanceName: tcpsentbytes.instance.{{ .Values.namespaces.adapter }}
      type: HISTOGRAM
      sample:
        expDecay:
          reservoirSize: 1024
          alpha: 0.015
    - name: tcpreceivedbytes
      instanceName: tcpreceivedbytes.instance.{{ .Values.namespaces.adapter }}
      type: HISTOGRAM
      sample:
        expDecay:
          reservoirSize: 1024
          alpha: 0.015
    - name: tcpconnectionsopened
      instanceName: tcpconnectionsopened.instance.{{ .Values.namespaces.adapter }}
      type: HISTOGRAM
      sample:
        expDecay:
          reservoirSize: 1024
          alpha: 0.015
    - name: tcpconnectionsclosed
      instanceName: tcpconnectionsclosed.instance.{{ .Values.namespaces.adapter }}
      type: HISTOGRAM
      sample:
        expDecay:
          reservoirSize: 1024
          alpha: 0.015
{{- end }}
