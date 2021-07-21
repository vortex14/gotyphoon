{{/*
Expand the name of the chart.
*/}}
{{- define "{{PROJECT_NAME}}.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "{{PROJECT_NAME}}.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "{{PROJECT_NAME}}.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "{{PROJECT_NAME}}.labels" -}}
helm.sh/chart: {{ include "{{PROJECT_NAME}}.chart" . }}
{{ include "{{PROJECT_NAME}}.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "{{PROJECT_NAME}}.selectorLabels" -}}
app.kubernetes.io/name: {{ include "{{PROJECT_NAME}}.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "{{PROJECT_NAME}}.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "{{PROJECT_NAME}}.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}



{{- define "imagePullSecret" -}}
{{- if .Values.hub.enabled }}
{{- printf "{\"auths\": {\"%s\": {\"auth\": \"%s\"}}}" .Values.hub.host (printf "%s:%s" .Values.hub.user .Values.hub.password | b64enc) | b64enc }}
{{- end }}
{{- end }}

