package short

import (
	"reflect"
	"strings"
)

// Short is a helper to return short, unique substrings when
// given a set of data to work with and the full value of
// the string to shorten. This can be useful to make indexes
// more human-friendly while still retaining their
// uniqueness and identifiability.
//
// A good example of where to use this library is with user-
// facing UUID's. It is often much easier to return a 6- or
// 7-character string and pass it around than it is to use
// the full 128-bit value.
type Short struct {
	// TODO: replace with a radix tree.
	data []string
}

// New creates a new shortener. It takes a slice of either
// strings or structs, and an optional field name. If using
// structs, the field name indicates which string field
// should be used.
func New(data interface{}, field string) *Short {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		panic("not a slice")
	}

	set := make([]string, v.Len())

	for i := 0; i < v.Len(); i++ {
		val := reflect.Indirect(v.Index(i))

		switch val.Kind() {
		case reflect.String:
			// If this is a []string we can short-circuit
			return &Short{data.([]string)}

		case reflect.Struct:
			// Accept structs

		default:
			panic("not a string or struct")
		}

		// Get the field value
		fieldVal := val.FieldByName(field)
		if !fieldVal.IsValid() {
			panic("missing field")
		}
		set[i] = fieldVal.String()
	}

	return &Short{set}
}

// min is the internal method used to retrieve the shortest
// possible string, given the length constraint.
//
// TODO: obviously needs massive optimizations
func (s *Short) min(in string, l int) string {
	var result string
OUTER:
	for i := 0; ; i++ {
		result += in[i*l : (i+1)*l]
		for _, item := range s.data {
			if strings.HasPrefix(item, result) {
				continue OUTER
			}
		}
		return result
	}
	return ""
}

// MinChunk is used to return the shortest substring in the
// chunk size provided. This means the minimum returned length
// is l, and the max is a multiple thereof. This is useful
// for keeping churn rate low with a frequently changing
// data set.
func (s *Short) MinChunk(in string, l int) string {
	return s.min(in, l)
}

// Min is used to return the shortest possible unique match
// from the data set.
func (s *Short) Min(in string) string {
	return s.min(in, 1)
}
