{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "csi-vault.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "csi-vault.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "chart.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}


{{- define "csi-vault.labels" -}}
app: "{{ template "csi-vault.name" . }}"
chart: "{{ .Chart.Name }}-{{ .Chart.Version }}"
release: {{ .Release.Name | quote}}
heritage: "{{ .Release.Service }}"
{{- end -}}

{{- define "csi-vault.attacher" -}}
{{- printf "%s-%s" (include "csi-vault.fullname" .) "attacher" | trunc 63 | trimSuffix "-" -}}
{{ end }}

{{- define "csi-vault.provisioner" -}}
{{- printf "%s-%s" (include "csi-vault.fullname" .) "provisioner" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define  "csi-vault.plugin" -}}
{{- printf "%s-%s" (include "csi-vault.fullname" .) "plugin" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "csi-vault.node" -}}
valueFrom:
  fieldRef:
    fieldPath: spec.nodeName
{{- end -}}