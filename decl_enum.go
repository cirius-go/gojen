// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package gojen

import (
	"fmt"
	"strings"
)

const (
	// StrategyTrunc is a Strategy of type trunc.
	StrategyTrunc Strategy = "trunc"
	// StrategyAppend is a Strategy of type append.
	StrategyAppend Strategy = "append"
	// StrategyAppendAtLast is a Strategy of type append_at_last.
	StrategyAppendAtLast Strategy = "append_at_last"
	// StrategyInit is a Strategy of type init.
	StrategyInit Strategy = "init"
)

var ErrInvalidStrategy = fmt.Errorf("not a valid Strategy, try [%s]", strings.Join(_StrategyNames, ", "))

var _StrategyNames = []string{
	string(StrategyTrunc),
	string(StrategyAppend),
	string(StrategyAppendAtLast),
	string(StrategyInit),
}

// StrategyNames returns a list of possible string values of Strategy.
func StrategyNames() []string {
	tmp := make([]string, len(_StrategyNames))
	copy(tmp, _StrategyNames)
	return tmp
}

// StrategyValues returns a list of the values for Strategy
func StrategyValues() []Strategy {
	return []Strategy{
		StrategyTrunc,
		StrategyAppend,
		StrategyAppendAtLast,
		StrategyInit,
	}
}

// String implements the Stringer interface.
func (x Strategy) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Strategy) IsValid() bool {
	_, err := ParseStrategy(string(x))
	return err == nil
}

var _StrategyValue = map[string]Strategy{
	"trunc":          StrategyTrunc,
	"append":         StrategyAppend,
	"append_at_last": StrategyAppendAtLast,
	"init":           StrategyInit,
}

// ParseStrategy attempts to convert a string to a Strategy.
func ParseStrategy(name string) (Strategy, error) {
	if x, ok := _StrategyValue[name]; ok {
		return x, nil
	}
	return Strategy(""), fmt.Errorf("%s is %w", name, ErrInvalidStrategy)
}

// MarshalText implements the text marshaller method.
func (x Strategy) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *Strategy) UnmarshalText(text []byte) error {
	tmp, err := ParseStrategy(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
