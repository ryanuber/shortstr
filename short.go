package short

import (
	"reflect"

	"github.com/armon/go-radix"
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
	tree *radix.Tree
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

	tree := radix.New()

	// Go over all of the data and insert our keys into
	// the tree.
	for i := 0; i < v.Len(); i++ {
		val := reflect.Indirect(v.Index(i))
		switch val.Kind() {
		case reflect.String:
			// No special handling required for strings

		case reflect.Struct:
			// If we have a struct, we need to attempt to
			// read the field value.
			val = val.FieldByName(field)
			if !val.IsValid() {
				panic("missing struct field")
			}

		default:
			panic("not a string or struct")
		}

		// Insert the value into the tree
		tree.Insert(val.String(), struct{}{})
	}

	return &Short{tree}
}

// min is the internal method used to retrieve the shortest
// possible string, given the length constraint.
func (s *Short) min(in string, l int) string {
	var result string
	for i := 0; ; i++ {
		// Add the next chunk of characters
		lidx := (i + 1) * l
		if lidx > len(in) {
			break
		}
		result += in[i*l : (i+1)*l]

		// Walk the tree. If anything is found by the given
		// result prefix, then the current result is ambiguous
		// and we need to add more characters.
		var ambiguous bool
		s.tree.WalkPrefix(result, func(string, interface{}) bool {
			ambiguous = true
			return true
		})
		if ambiguous {
			continue
		}

		// We got an unambiguous result, so return it
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
// from the data set. If an empty string is returned, then
func (s *Short) Min(in string) string {
	return s.min(in, 1)
}
