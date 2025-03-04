package decl

import "github.com/cirius-go/gojen"

var DTODecl = &gojen.D{
	Path:    "internal/dto/{{ sSnake .Domain }}.go",
	Name:    "dto",
	Require: []string{"Domain"},
	Templates: []*gojen.T{
		{
			Name:     "initDtoFile",
			Template: `package dto`,
			Strategy: gojen.StrategyInit,
		},
		{
			Name:    "createDto",
			Require: []string{"Domain", "Methods"},
			Template: `{{ $domain := .Domain | sIniCamel -}}
{{- range .Methods -}}
  {{ $method := . | sIniCamel }}
  {{- $req := printf "%s%sReq" $method $domain -}}
  {{- $res := printf "%s%sRes" $method $domain }}
type (
  // {{ $req }} is the request data of {{ $domain }}.{{ $method }}.
  {{ $req }} struct {}

  // {{ $res }} is the response data of {{ $domain }}.{{ $method }}.
  {{ if eq $method "List" -}}
  {{ $res }} = ListingRes[List{{ $domain }}Item]

  // List{{ $domain }}Item is the item for list response of {{ $domain }}.
  List{{ $domain }}Item struct {}
  {{- else if or (eq $method "Update") (eq $method "Create") -}}
    {{ $res }} = Get{{ $domain }}Res
  {{- else -}}
    {{ $res }} struct {}
  {{- end }}
)
{{ end -}}`,
			Strategy: gojen.StrategyAppendAtPos,
			Output: map[string]*gojen.Output{
				"createView": {
					Path:     "internal/repo/model/{{ sSnake .Domain }}.go",
					Template: "",
				},
			},
		},
	},
	Description: "DTO declaration",
}
