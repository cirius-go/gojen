package gojen

import "gopkg.in/yaml.v2"

type (
	// BuildState contains the build state of a seq.
	State struct {
		seq *Seq `yaml:"-"`
		d   *D   `yaml:"-"`
		e   *E   `yaml:"-"`

		Strategy      Strategy `yaml:"strategy"`
		DName         string   `yaml:"d_name"`
		EName         string   `yaml:"e_name"`
		RawEAlias     string   `yaml:"e_alias"`
		ParsedEAlias  string   `yaml:"parsed_e_alias"`
		Args          Args     `yaml:"args"`
		ForwardedArgs Args     `yaml:"forwarded_args"`
		RawPath       string   `yaml:"raw_path"`
		RawTmpl       string   `yaml:"raw_tmpl"`
		ParsedPath    string   `yaml:"parsed_path"`
		ParsedTmpl    string   `yaml:"parsed_tmpl"`
	}
)

func (s *State) String() string {
	v, err := yaml.Marshal(s)
	if err != nil {
		return ""
	}

	return string(v)
}
