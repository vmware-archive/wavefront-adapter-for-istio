{{/* Generate http metrics */}}
{{- define "metrics.http" }}
    - name: requestsize
      instanceName: requestsize.instance.istio-system
      type: HISTOGRAM
      sample:
        expDecay:
          reservoirSize: 1024
          alpha: 0.015
    - name: requestcount
      instanceName: requestcount.instance.istio-system
      type: DELTA_COUNTER
    - name: requestduration
      instanceName: requestduration.instance.istio-system
      type: HISTOGRAM
      sample:
        expDecay:
          reservoirSize: 1024
          alpha: 0.015
    - name: responsesize
      instanceName: responsesize.instance.istio-system
      type: HISTOGRAM
      sample:
        expDecay:
          reservoirSize: 1024
          alpha: 0.015
{{- end }}
