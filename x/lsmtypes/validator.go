package lsmtypes

import "gopkg.in/yaml.v2"

// String implements the Stringer interface for a Validator object.
func (v Validator) String() string {
	out, _ := yaml.Marshal(v)
	return string(out)
}
