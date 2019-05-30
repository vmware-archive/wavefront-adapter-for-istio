{{/* Generate http metrics */}}
{{- define "metrics.http" }}
    - name: requestsize
      instanceName: requestsize.instance.{{ .Values.adapterNamespace }}
      type: HISTOGRAM
      sample:
        expDecay:
          reservoirSize: 1024
          alpha: 0.015
    - name: requestcount
      instanceName: requestcount.instance.{{ .Values.adapterNamespace }}
      type: DELTA_COUNTER
    - name: requestduration
      instanceName: requestduration.instance.{{ .Values.adapterNamespace }}
      type: HISTOGRAM
      sample:
        expDecay:
          reservoirSize: 1024
          alpha: 0.015
    - name: responsesize
      instanceName: responsesize.instance.{{ .Values.adapterNamespace }}
      type: HISTOGRAM
      sample:
        expDecay:
          reservoirSize: 1024
          alpha: 0.015
{{- end }}
