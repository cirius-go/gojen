package decl

import "github.com/cirius-go/gojen"

var (
	ServiceDecl = &gojen.D{
		Path:    "internal/service/{{ sSnake .Domain }}.go",
		Name:    "svc",
		Require: []string{"Domain"},
		Templates: []*gojen.T{
			{
				Name: "initSvcFile",
				Template: `package service

// {{ sIniCamel .Domain }} service.
type {{ sIniCamel .Domain }} struct {
  Service
  uow uow.UnitOfWork
}

// New{{ sIniCamel .Domain }} creates a new {{ .Domain }}.
func New{{ sIniCamel .Domain }}(uow uow.UnitOfWork, enf Enforcer) *{{ sIniCamel .Domain }} {
  s := &{{ sIniCamel .Domain }}{
	Service: NewServiceRBAC(enf, model.Object{{ sIniCamel .Domain }}),
    uow: uow,
  }

  s.AssignRole(model.RoleAdmin, )

  return s
}
`,
				Strategy: gojen.StrategyInit,
				Output: map[string]*gojen.Output{
					"newObject": {
						Path: "internal/repo/model/rbac.go",
						Template: `
{{- $domain := .Domain | sIniCamel -}}
Object{{ $domain }} Object = "{{ $domain }}"`,
					},
				},
			},
			{
				Name:    "createSvcHandler",
				Require: []string{"Domain", "Methods"},
				Template: `{{- $domain := .Domain | sIniCamel -}}
{{ $method := .Method | sIniCamel }}
func (s *{{ $domain }}) {{ $method }}(ctx context.Context, req *dto.{{ $method }}{{ $domain }}Req) (*dto.{{ $method }}{{ $domain }}Res, error) {
  if err := s.Enforce(ctx, model.Action{{ $method }}{{ $domain }}); err != nil {
		return nil, err
  }

  if err := s.Validate(ctx, req); err != nil {
    return nil, err
  }
  panic("not implemented")
}
{{- end -}}`,
				Strategy: gojen.StrategyAppendAtPos,
				Output: map[string]*gojen.Output{
					"newAction": {
						Path: "internal/repo/model/rbac.go",
						Template: `
{{- $domain := .Domain | sIniCamel -}}
{{- $method := .Method | sIniCamel -}}
Action{{ $method }}{{ $domain }} Action = "{{ $method }}{{ $domain }}"`,
					},
				},
			},
		},
		Description: "Service declaration",
	}
)
