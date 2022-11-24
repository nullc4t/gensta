package templates

var (
	StructTemplate = `
// {{ .StructName }} is an exchange struct
type {{ .StructName }} struct {
{{- range .Fields }}
	{{ .Name }} {{ .Type.String }}
{{- else -}}
{{- end -}}
}
`
	ProtocolTemplate = `package {{ .Package }}

{{ range .Structs }}
{{ template "struct" . }}
{{ if .Fields }}

// New{{ .StructName }} is a constructor for {{ .StructName }}
func New{{ .StructName }} ({{ struct_constructor_args .Fields }}) {{ .StructName }} {
	return {{ .StructName }}{ {{- struct_constructor_return .Fields -}} }
}

// Args is a shortcut method that returns args to original interface's method
func (r {{ .StructName }}) Args(ctx context.Context) (context.Context, {{ struct_return_types .Fields }}) {
	return ctx, {{ struct_return_args .Fields }}
}

{{ end }}
{{ end }}
`
)
