package global

{{- if .HasGlobal }}

import "zc-admin/server/plugin/{{ .Snake}}/config"

var GlobalConfig = new(config.{{ .PlugName}})
{{ end -}}
