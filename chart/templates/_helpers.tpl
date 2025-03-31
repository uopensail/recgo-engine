{{/*
Expand the name of the chart.
*/}}
{{- define "recgo-engine.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "recgo-engine.fullname" -}}
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
{{- define "recgo-engine.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "recgo-engine.labels" -}}
helm.sh/chart: {{ include "recgo-engine.chart" . }}
{{ include "recgo-engine.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "recgo-engine.selectorLabels" -}}
app.kubernetes.io/name: {{ include "recgo-engine.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "recgo-engine.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "recgo-engine.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/* PVC existing, emptyDir, Dynamic */}}
{{- define "recgo-engine.pvc" -}}
{{- if or (not .Values.recgo-engine.persistence.enabled) (eq .Values.recgo-engine.persistence.type "emptyDir") -}}
          emptyDir: {}
{{- else if and .Values.recgo-engine.persistence.enabled .Values.recgo-engine.persistence.existingClaim -}}
          persistentVolumeClaim:
            claimName: {{ .Values.recgo-engine.persistence.existingClaim }}
{{- else if and .Values.recgo-engine.persistence.enabled (eq .Values.recgo-engine.persistence.type "dynamic")  -}}
          persistentVolumeClaim:
            claimName: {{ include "recgo-engine.fullname" . }}
{{- end }}
{{- end }}

{{- define "recgo-engine.mongo.service" -}}
{{ include "recgo-engine.fullname" . }}-{{- if eq .Values.mongodb.architecture "replicaset" -}}mongodb-headless{{- else -}}mongodb{{- end -}}
{{- end }}

{{- define "recgo-engine.mongo.uri" -}}
{{- if .Values.mongodb.enabled -}}
{{- if ne .Values.mongodb.architecture "replicaset" -}}
{{- printf "mongodb://%s:%s@%s:27017/%s" .Values.mongodb.auth.username .Values.mongodb.auth.password (include "recgo-engine.mongo.service" .) .Values.mongodb.auth.database }}
{{- else }}
{{- printf "mongodb://%s:%s@%s:27017/%s?replicaSet=%s" .Values.mongodb.auth.username .Values.mongodb.auth.password (include "recgo-engine.mongo.service" .) .Values.mongodb.auth.database .Values.mongodb.replicaSetName }}
{{- end }}
{{- else }}
{{- .Values.recgo-engine.externalMongodbUri }}
{{- end }}
{{- end }}

{{- define "recgo-engine.mongo.env" -}}
- name: MONGODB_PORT
  value: '27017'
- name: MONGODB_HOST
  value: {{ include "recgo-engine.mongo.service" . }}
{{- end }}
