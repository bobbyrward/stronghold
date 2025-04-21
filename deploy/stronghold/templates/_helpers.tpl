{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "stronghold.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "stronghold.labels" -}}
app.kubernetes.io/name: "stronghold"
app.kubernetes.io/instance: {{ .Release.Name }}
helm.sh/chart: {{ include "stronghold.chart" . }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "stronghold.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "stronghold.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "stronghold.image" -}}
{{ .Values.images.stronghold.repository }}:{{ .Values.images.stronghold.tag | default .Chart.AppVersion }}
{{- end }}
