{{/* Generate http metrics */}}
{{- define "metrics.http" }}
    - name: requestsize
      instanceName: requestsize.instance.{{ .Values.namespaces.adapter }}
      type: HISTOGRAM
      sample:
        expDecay:
          reservoirSize: 1024
          alpha: 0.015
    - name: requestcount
      instanceName: requestcount.instance.{{ .Values.namespaces.adapter }}
      type: DELTA_COUNTER
    - name: requestduration
      instanceName: requestduration.instance.{{ .Values.namespaces.adapter }}
      type: HISTOGRAM
      sample:
        expDecay:
          reservoirSize: 1024
          alpha: 0.015
    - name: responsesize
      instanceName: responsesize.instance.{{ .Values.namespaces.adapter }}
      type: HISTOGRAM
      sample:
        expDecay:
          reservoirSize: 1024
          alpha: 0.015
{{- end }}
