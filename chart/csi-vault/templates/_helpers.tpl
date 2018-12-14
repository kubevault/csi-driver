{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "vault.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "vault.fullname" -}}
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


{{- define "vault.labels" -}}
chart: "{{ .Chart.Name }}-{{ .Chart.Version }}"
release: {{ .Release.Name | quote}}
heritage: "{{ .Release.Service }}"
{{- end -}}

{{- define "vault.attacher" -}}
{{- printf "%s-%s" .Release.Name "attacher" | trunc 63 | trimSuffix "-" -}}
{{ end }}

{{- define "vault.provisioner" -}}
{{- printf "%s-%s" .Release.Name "provisioner" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define  "vault.plugin" -}}
{{- printf "%s-%s" .Release.Name "plugin" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "vault.node" -}}
valueFrom:
  fieldRef:
    fieldPath: spec.nodeName
{{- end -}}