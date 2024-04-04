package types

import (
	"strconv"
	"strings"
)

type FieldValue struct {
	Field  string
	Value  string
	Negate bool
}

func (e Event) ResolveAllFieldValues(fvs []FieldValue) bool {
	for _, fv := range fvs {
		if !e.resolveFieldValue(fv) {
			return false
		}
	}
	return true
}

func (e Event) ResolveAnyFieldValues(fvs []FieldValue) bool {
	for _, fv := range fvs {
		if e.resolveFieldValue(fv) {
			return true
		}
	}
	return false
}

func (e Event) resolveFieldValue(fv FieldValue) bool {

	if strings.ToLower(fv.Field) == "eventtype" {
		v, err := strconv.ParseInt(fv.Value, 10, 32)
		if err != nil {
			return fv.Negate
		}
		if v == int64(e.EventType) {
			return !fv.Negate
		}
		return fv.Negate
	}

	if strings.ToLower(fv.Field) == "eventstatus" {
		v, err := strconv.ParseInt(fv.Value, 10, 32)
		if err != nil {
			return fv.Negate
		}
		if v == int64(e.EventType) {
			return !fv.Negate
		}
		return fv.Negate
	}

	if strings.ToLower(fv.Field) == "module" && fv.Value == e.Module {
		return !fv.Negate
	}
	if strings.ToLower(fv.Field) == "identifier" && fv.Value == e.Identifier {
		return !fv.Negate
	}
	if strings.ToLower(fv.Field) == "chainid" && fv.Value == e.ChainId {
		return !fv.Negate
	}

	return fv.Negate
}
