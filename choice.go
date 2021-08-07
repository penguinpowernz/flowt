package flowt

import (
	"encoding/json"
	"strings"

	"github.com/Jeffail/gabs"
)

type Choice struct {
	And []Choice
	Not *Choice

	Variable string

	StringEquals *string

	StringEqualsPath       string
	NumericGreaterThanPath string
	NumericLessThanPath    string

	NumericGreaterThanEquals float64
	NumericLessThan          float64

	IsNull    *bool
	IsString  *bool
	IsBoolean *bool
	IsNumeric *bool
	IsPresent *bool

	Next string
}

func (ch Choice) IsSatisfied(input *gabs.Container) bool {
	switch {
	case len(ch.And) > 0:
		for _, ch := range ch.And {
			if !ch.IsSatisfied(input) {
				return false
			}
		}
		return true

	case ch.Not != nil:
		return !ch.Not.IsSatisfied(input)

	case ch.Variable != "":
		varpath := strings.TrimPrefix(ch.Variable, "$.")
		switch {
		case ch.StringEqualsPath != "":
			sv, ok := input.Path(varpath).Data().(string)
			if !ok {
				return false
			}

			sexp, ok := input.Path(varpath).Data().(string)
			if !ok {
				return false
			}

			if sv == sexp {
				return true
			}

		case ch.StringEquals != nil:
			s, ok := input.Path(varpath).Data().(string)
			if !ok {
				return false
			}
			if s == *ch.StringEquals {
				return true
			}

		case ch.IsNull != nil:
			if !input.ExistsP(varpath) {
				return false
			}

			if input.Path(varpath).Data() == nil {
				return true
			}

		case ch.IsPresent != nil:
			if input.ExistsP(varpath) == *ch.IsPresent {
				return true
			}

		case ch.IsNumeric != nil:
			if testIsNumeric(input.Path(varpath).Data()) == *ch.IsNumeric {
				return true
			}
		}

	}

	return false
}

func testIsNumeric(x interface{}) bool {
	switch x.(type) {
	case int, int64, int32, float64, json.Number, float32:
		return true
	}
	return false
}
