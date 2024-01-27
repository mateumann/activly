{{/*
Expand the name of the chart.
*/}}
{{- define "activly.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "activly.fullname" -}}
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
{{- define "activly.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "activly.labels" -}}
helm.sh/chart: {{ include "activly.chart" . }}
{{ include "activly.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "activly.selectorLabels" -}}
app.kubernetes.io/name: {{ include "activly.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "activly.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "activly.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}


{{/*
Return the proper image name (registry/repository:tag)
{{ include "activly.images.image" ( dict "imageRoot" .Values.path.to.the.image "root" $ ) }}
*/}}
{{- define "activly.images.image" -}}
{{- $registryName := .root.Values.global.image.registry -}}
{{- $repositoryName := .imageRoot.repository -}}
{{- $separator := ":" -}}
{{- $termination := default .root.Chart.AppVersion .imageRoot.tag | toString -}}
{{- if .imageRoot.registry }}
    {{- $registryName = .imageRoot.registry -}}
{{- end -}}
{{- if .imageRoot.digest }}
    {{- $separator = "@" -}}
    {{- $termination = .imageRoot.digest | toString -}}
{{- end -}}
{{- printf "%s/%s%s%s" $registryName $repositoryName $separator $termination -}}
{{- end -}}

{{/*
Return env vars block with postgres connection parameters.
{{ include "activly.postgres.env" ( dict "root" $ ) }}
*/}}
{{- define "activly.postgres.env" -}}
- name: POSTGRES_HOST
  value: {{ (.root).Values.postgresql.fullnameOverride }}
- name: "POSTGRES_PORT"
  value: {{ (.root).Values.postgresql.primary.service.ports.postgresql | quote }}
- name: POSTGRES_DB
  value: {{ (.root).Values.postgresql.auth.database }}
- name: POSTGRES_USER
  value: {{ (.root).Values.postgresql.auth.username }}
- name: POSTGRES_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ (.root).Values.postgresql.auth.existingSecret }}
      key: {{ (.root).Values.postgresql.auth.secretKeys.userPasswordKey }}
{{- end -}}