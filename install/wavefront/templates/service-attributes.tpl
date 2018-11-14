{{/* Generate attributes for the source and destination service */}}
{{- define "attributes.service" }}
      reporter: conditional((context.reporter.kind | "inbound") == "outbound", "client", "server")
      source_service: source.workload.name | "unknown"
      source_service_namespace: source.workload.namespace | "unknown"
      source_version: source.labels["version"] | "unknown"
      destination_service: destination.service.name | "unknown"
      destination_service_namespace: destination.service.namespace | "unknown"
      destination_version: destination.labels["version"] | "unknown"
{{- end }}
