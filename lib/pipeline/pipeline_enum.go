// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package pipeline

import (
	"fmt"
	"strings"
)

const (
	// PluralizeTypePlural is a PluralizeType of type plural.
	PluralizeTypePlural PluralizeType = "plural"
	// PluralizeTypeSingular is a PluralizeType of type singular.
	PluralizeTypeSingular PluralizeType = "singular"
	// PluralizeTypeIrregular is a PluralizeType of type irregular.
	PluralizeTypeIrregular PluralizeType = "irregular"
)

var ErrInvalidPluralizeType = fmt.Errorf("not a valid PluralizeType, try [%s]", strings.Join(_PluralizeTypeNames, ", "))

var _PluralizeTypeNames = []string{
	string(PluralizeTypePlural),
	string(PluralizeTypeSingular),
	string(PluralizeTypeIrregular),
}

// PluralizeTypeNames returns a list of possible string values of PluralizeType.
func PluralizeTypeNames() []string {
	tmp := make([]string, len(_PluralizeTypeNames))
	copy(tmp, _PluralizeTypeNames)
	return tmp
}

// PluralizeTypeValues returns a list of the values for PluralizeType
func PluralizeTypeValues() []PluralizeType {
	return []PluralizeType{
		PluralizeTypePlural,
		PluralizeTypeSingular,
		PluralizeTypeIrregular,
	}
}

// String implements the Stringer interface.
func (x PluralizeType) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x PluralizeType) IsValid() bool {
	_, err := ParsePluralizeType(string(x))
	return err == nil
}

var _PluralizeTypeValue = map[string]PluralizeType{
	"plural":    PluralizeTypePlural,
	"singular":  PluralizeTypeSingular,
	"irregular": PluralizeTypeIrregular,
}

// ParsePluralizeType attempts to convert a string to a PluralizeType.
func ParsePluralizeType(name string) (PluralizeType, error) {
	if x, ok := _PluralizeTypeValue[name]; ok {
		return x, nil
	}
	return PluralizeType(""), fmt.Errorf("%s is %w", name, ErrInvalidPluralizeType)
}

// MarshalText implements the text marshaller method.
func (x PluralizeType) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *PluralizeType) UnmarshalText(text []byte) error {
	tmp, err := ParsePluralizeType(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
