{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "stolon.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "stolon.fullname" -}}
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
{{- define "stolon.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}


{{/* Create the name of the service account to use */}}
{{- define "stolon.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "stolon.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{- define "stolon.clusterName" -}}
{{- if .Values.clusterName -}}
    {{- .Values.clusterName | trunc 63 | trimSuffix "-" -}}
{{- else -}}
    {{- template "stolon.fullname" . -}}
{{- end -}}
{{- end -}}

{{- define "stolon.clusterSpec" -}}
{{- range $key, $value := .Values.clusterSpec }} {{ $key | quote }}: {{- $tp := typeOf $value }} {{- if eq $tp "string" }} {{ $value | quote }} {{- else if eq $tp "int" }} {{ $value | int64 }} {{- else if eq $tp "float64" }} {{ $value | int64 }} {{- else if eq $tp "[]interface {}" }} {{- $numOut := len $value }} {{- $numOut := sub $numOut 1 }} [{{- range $b, $val := $value }} {{- $i := int64 $b }} {{- if eq $i $numOut }} {{ $val | quote }} {{- else }} {{ $val | quote }}, {{- end }} {{- end }}] {{- else }} {{ $value }} {{- end }}, {{- end }} "pgParameters": {{ toJson .Values.pgParameters }}
{{- end -}}

{{/* Telegraf helpers. Original https://github.com/helm/charts/blob/master/stable/telegraf/templates/_helpers.tpl */}}

{{- define "telegraf_global_tags" -}}
{{- if . -}}
[global_tags]
  {{- range $key, $val := . }}
      {{ $key }} = {{ $val | quote }}
  {{- end }}
{{- end }}
{{- end -}}

{{- define "telegraf_agent" -}}
[agent]
{{- range $key, $value := . -}}
  {{- $tp := typeOf $value }}
  {{- if eq $tp "string"}}
      {{ $key }} = {{ $value | quote }}
  {{- end }}
  {{- if eq $tp "float64"}}
      {{ $key }} = {{ $value | int64 }}
  {{- end }}
  {{- if eq $tp "int"}}
      {{ $key }} = {{ $value | int64 }}
  {{- end }}
  {{- if eq $tp "bool"}}
      {{ $key }} = {{ $value }}
  {{- end }}
{{- end }}
{{- end -}}

{{- define "telegraf_outputs" -}}
{{- range $outputIdx, $configObject := . -}}
{{- range $output, $config := . }}
    [[outputs.{{ $output }}]]
  {{- if $config }}
  {{- $tp := typeOf $config -}}
  {{- if eq $tp "map[string]interface {}" -}}
    {{- range $key, $value := $config -}}
      {{- $tp := typeOf $value }}
      {{- if eq $tp "string"}}
      {{ $key }} = {{ $value | quote }}
      {{- end }}
      {{- if eq $tp "float64"}}
      {{ $key }} = {{ $value | int64 }}
      {{- end }}
      {{- if eq $tp "int"}}
      {{ $key }} = {{ $value | int64 }}
      {{- end }}
      {{- if eq $tp "bool"}}
      {{ $key }} = {{ $value }}
      {{- end }}
      {{- if eq $tp "[]interface {}" }}
      {{ $key }} = [
          {{- $numOut := len $value }}
          {{- $numOut := sub $numOut 1 }}
          {{- range $b, $val := $value }}
            {{- $i := int64 $b }}
            {{- if eq $i $numOut }}
        {{ $val | quote }}
            {{- else }}
        {{ $val | quote }},
            {{- end }}
          {{- end }}
      ]
      {{- end }}
    {{- end }}
  {{- end }}
  {{- end }}
  {{- end }}
{{- end }}
{{- end -}}

{{- define "telegraf_inputs" -}}
{{- range $inputIdx, $configObject := . -}}
    {{- range $input, $config := . -}}

    [[inputs.{{- $input }}]]
    {{- if $config -}}
    {{- $tp := typeOf $config -}}
    {{- if eq $tp "map[string]interface {}" -}}
        {{- range $key, $value := $config -}}
          {{- $tp := typeOf $value -}}
          {{- if eq $tp "string" }}
      {{ $key }} = {{ $value | quote }}
          {{- end }}
          {{- if eq $tp "float64" }}
      {{ $key }} = {{ $value | int64 }}
          {{- end }}
          {{- if eq $tp "int" }}
      {{ $key }} = {{ $value | int64 }}
          {{- end }}
          {{- if eq $tp "bool" }}
      {{ $key }} = {{ $value }}
          {{- end }}
          {{- if eq $tp "[]interface {}" }}
      {{ $key }} = [
              {{- $numOut := len $value }}
              {{- $numOut := sub $numOut 1 }}
              {{- range $b, $val := $value }}
                {{- $i := int64 $b }}
                {{- $tp := typeOf $val }}
                {{- if eq $i $numOut }}
                  {{- if eq $tp "string" }}
        {{ $val | quote }}
                  {{- end }}
                  {{- if eq $tp "float64" }}
        {{ $val | int64 }}
                  {{- end }}
                {{- else }}
                  {{- if eq $tp "string" }}
        {{ $val | quote }},
                  {{- end}}
                  {{- if eq $tp "float64" }}
        {{ $val | int64 }},
                  {{- end }}
                {{- end }}
              {{- end }}
      ]
          {{- end }}
          {{- if eq $tp "map[string]interface {}" }}
      [[inputs.{{ $input }}.{{ $key }}]]
            {{- range $k, $v := $value }}
              {{- $tps := typeOf $v }}
              {{- if eq $tps "string" }}
        {{ $k }} = {{ $v }}
              {{- end }}
              {{- if eq $tps "[]interface {}"}}
        {{ $k }} = [
                {{- $numOut := len $value }}
                {{- $numOut := sub $numOut 1 }}
                {{- range $b, $val := $v }}
                  {{- $i := int64 $b }}
                  {{- if eq $i $numOut }}
            {{ $val | quote }}
                  {{- else }}
            {{ $val | quote }},
                  {{- end }}
                {{- end }}
        ]
              {{- end }}
              {{- if eq $tps "map[string]interface {}"}}
        [[inputs.{{ $input }}.{{ $key }}.{{ $k }}]]
                {{- range $foo, $bar := $v }}
            {{ $foo }} = {{ $bar | quote }}
                {{- end }}
              {{- end }}
            {{- end }}
          {{- end }}
          {{- end }}
    {{- end }}
    {{- end }}
    {{ end }}
{{- end }}
{{- end -}}

{{- define "telegraf_processors" -}}
{{- range $processorIdx, $configObject := . -}}
    {{- range $processor, $config := . -}}

    [[processors.{{- $processor }}]]
    {{- if $config -}}
    {{- $tp := typeOf $config -}}
    {{- if eq $tp "map[string]interface {}" -}}
        {{- range $key, $value := $config -}}
          {{- $tp := typeOf $value -}}
          {{- if eq $tp "string" }}
      {{ $key }} = {{ $value | quote }}
          {{- end }}
          {{- if eq $tp "float64" }}
      {{ $key }} = {{ $value | int64 }}
          {{- end }}
          {{- if eq $tp "int" }}
      {{ $key }} = {{ $value | int64 }}
          {{- end }}
          {{- if eq $tp "bool" }}
      {{ $key }} = {{ $value }}
          {{- end }}
          {{- if eq $tp "[]interface {}" }}
      {{ $key }} = [
              {{- $numOut := len $value }}
              {{- $numOut := sub $numOut 1 }}
              {{- range $b, $val := $value }}
                {{- $i := int64 $b }}
                {{- $tp := typeOf $val }}
                {{- if eq $i $numOut }}
                  {{- if eq $tp "string" }}
        {{ $val | quote }}
                  {{- end }}
                  {{- if eq $tp "float64" }}
        {{ $val | int64 }}
                  {{- end }}
                {{- else }}
                  {{- if eq $tp "string" }}
        {{ $val | quote }},
                  {{- end}}
                  {{- if eq $tp "float64" }}
        {{ $val | int64 }},
                  {{- end }}
                {{- end }}
              {{- end }}
      ]
          {{- end }}
          {{- if eq $tp "map[string]interface {}" }}
      [[processors.{{ $processor }}.{{ $key }}]]
            {{- range $k, $v := $value }}
              {{- $tps := typeOf $v }}
              {{- if eq $tps "string" }}
        {{ $k }} = {{ $v | quote }}
              {{- end }}
              {{- if eq $tps "[]interface {}"}}
        {{ $k }} = [
                {{- $numOut := len $value }}
                {{- $numOut := sub $numOut 1 }}
                {{- range $b, $val := $v }}
                  {{- $i := int64 $b }}
                  {{- if eq $i $numOut }}
            {{ $val | quote }}
                  {{- else }}
            {{ $val | quote }},
                  {{- end }}
                {{- end }}
        ]
              {{- end }}
              {{- if eq $tps "map[string]interface {}"}}
        [processors.{{ $processor }}.{{ $key }}.{{ $k }}]
                {{- range $foo, $bar := $v }}
                {{- $tp := typeOf $bar -}}
                {{- if eq $tp "string" }}
            {{ $foo }} = {{ $bar | quote }}
                {{- end }}
                {{- if eq $tp "int" }}
            {{ $foo }} = {{ $bar }}
                {{- end }}
                {{- if eq $tp "float64" }}
            {{ $foo }} = {{ int64 $bar }}
                {{- end }}
                {{- end }}
              {{- end }}
            {{- end }}
          {{- end }}
          {{- end }}
    {{- end }}
    {{- end }}
    {{ end }}
{{- end }}
{{- end -}}
