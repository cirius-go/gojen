package decl

import "github.com/cirius-go/gojen"

var APIDecl = &gojen.D{
	Path:        "internal/api/{{ gopkg .BaseAPI }}api/{{ sSnake .Domain }}.go",
	Description: "API declaration",
	Name:        "api",
	Require:     []string{"Domain", "BaseAPI"},
	Templates: []*gojen.T{
		{
			Name: "initApiFile",
			Template: `
package {{ gopkg .BaseAPI }}api

import (
	"github.com/labstack/echo/v4"
)

// {{ sIniCamel .Domain }} is a handler struct that manages HTTP,... requests
// related to {{ .Domain }}.
type {{ sIniCamel .Domain }} struct {
  api.API //base for all API handlers
  svc {{ sIniCamel .Domain }}Service
}


// New{{ sIniCamel .Domain }} creates a new {{ .Domain }}.
func New{{ sIniCamel .Domain }}(svc {{ sIniCamel .Domain }}Service) *{{ sIniCamel .Domain }} {
  return &{{ sIniCamel .Domain }}{
    svc: svc,
  }
}

// RegisterHTTP registers {{ sIniCamel .Domain }} handlers to echo group.
func (h *{{ sIniCamel .Domain }}) RegisterHTTP(g *echo.Group) {
  g = g.Group("/{{ .Domain | pSnake }}")

	// +gojen:input=api.createHandler->echoBinding
}
`,
			Strategy: gojen.StrategyInit,
		},
		{
			Name:    "createApiHandler",
			Require: []string{"Domain", "Method", "HTTPMethod", "BaseAPI", "Slug"},
			Template: `{{- $domain := .Domain | sIniCamel -}}
{{- $baseAPI := .BaseAPI | gopkg -}}
{{ $method := .Method | sIniCamel }}
{{ $httpMethod := .HTTPMethod | upper }}
{{ $slug := .Slug }}

// {{ $method }}
//
//	@id {{ $baseAPI }}-{{ $domain | pKebab }}-{{ $method | sKebab | lower }}
//	@Summary {{ $method }}
//	@Description {{ $method }}
//	@Tags {{ $baseAPI }}/{{ $domain | pKebab | lower }}
//	@Accept json
//	@Produce json
//	@Security BearerAuth
{{- if or (eq $httpMethod "POST") (eq $httpMethod "DELETE") (eq $httpMethod "PUT") (eq $httpMethod "PATCH") }}
//	@Param id path string true "ID"
//	@Param Payload body dto.{{ $method }}{{ $domain }}Req true "JSON Request Payload"
{{- end }}
{{- if (eq $httpMethod "GET") }}
//	@param Payload query dto.{{ $method }}{{ $domain }}Req true "Request Payload"
{{- end }}
//	@Success 200 {object} dto.{{ $method }}{{ $domain }}Res "JSON Response Payload"
//	@Failure 400 {object} dto.ErrorRes "JSON Response Payload"
//	@Failure 500 {object} dto.ErrorRes "JSON Response Payload"
//	@Router /{{ $baseAPI }}/{{ $domain | pKebab | lower }}{{ $slug | parseEchoSlug }} [{{ $httpMethod }}]
func (h *{{ $domain }}) {{ $method }}(c echo.Context) error {
  return api.JSONHandlerFunc(h.svc.{{ $method }})(c)
}`,
			Output: map[string]*gojen.Output{
				"echoBinding": {
					Path: "internal/api/{{ gopkg .BaseAPI }}api/{{ sSnake .Domain }}.go",
					Template: `
{{- $method := .Method | sIniCamel }}
{{- $httpMethod := .HTTPMethod | upper }}
{{- $slug := .Slug -}}
g.{{ $httpMethod }}("{{ $slug }}", h.{{ $method }})`,
				},
				"serviceHandler": {
					Path: "internal/api/{{ gopkg .BaseAPI }}api/interface.go",
					Template: `
{{- $domain := .Domain | sIniCamel -}}
{{- $method := .Method | sIniCamel -}}
{{$method }}(ctx context.Context, req *dto.{{ $method }}{{ $domain }}Req) (*dto.{{ $method }}{{ $domain }}Res, error)`,
				},
			},
			Strategy: gojen.StrategyAppendAtPos,
		},
		{
			Path:     "internal/api/{{ gopkg .BaseAPI }}api/interface.go",
			Name:     "initIntfFile",
			Template: `package {{ gopkg .BaseAPI }}api`,
			Strategy: gojen.StrategyInit,
		},
		{
			Path:  "internal/api/{{ gopkg .BaseAPI }}api/interface.go",
			Name:  "createIntf",
			Alias: "{{ sIniCamel .Domain }}_decl",
			Template: `
// {{ sIniCamel .Domain }}Service contains required methods to handle api requests.
type {{ sIniCamel .Domain }}Service interface {
	// +gojen:input=api.createHandler->serviceHandler
}`,
			Strategy: gojen.StrategyAppendAtPos,
		},
	},
}
