package shortstr

import (
	"fmt"
	"reflect"

	"github.com/armon/go-radix"
)

// Shortener is a helper to return short, unique substrings when given a set of
// data to work with and the full value of the string to shorten. This can be
// useful to make indexes more human-friendly while still retaining their
// uniqueness and identifiability.
//
// A good example of where to use this library is with user-facing UUID's. It
// is often much easier to return a 6- or 7-character string and pass it around
// than it is to use the full 128-bit value.
type Shortener struct {
	tree     *radix.Tree
	field    string
	dataType reflect.Type
}

// New creates a new Shortener given a set of structs and a field name to use
// for comparison. If the input is not a slice of structs (or struct pointers),
// or the specified field does not exist in the struct, New will panic.
func New(data interface{}, field string) *Shortener {
	// Check that we have a slice
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		panic("not a slice")
	}
	elem := v.Type().Elem()

	var fieldVal reflect.StructField
	var ok bool

	// Check the slice type
	switch elem.Kind() {
	case reflect.Struct:
		// Struct values are allowed
		fieldVal, ok = elem.FieldByName(field)

	case reflect.Ptr:
		// Struct pointers are allowed
		if e := elem.Elem(); e.Kind() == reflect.Struct {
			fieldVal, ok = e.FieldByName(field)
			break
		}
		fallthrough

	default:
		panic("not a struct slice")
	}

	// Make sure our structs actually have the field
	if !ok {
		panic("invalid struct field")
	}
	if fieldVal.Type.Kind() != reflect.String {
		panic("struct field must be type string")
	}

	// Create the shortener
	s := &Shortener{
		tree:     radix.New(),
		field:    field,
		dataType: elem,
	}

	// Go over all of the data and insert our keys into the tree.
	for i := 0; i < v.Len(); i++ {
		s.Insert(v.Index(i))
	}

	return s
}

// NewStrings creates a new Shortener from a string slice.
func NewStrings(data []string) *Shortener {
	tree := radix.New()
	for _, s := range data {
		tree.Insert(s, struct{}{})
	}
	return &Shortener{
		tree:     tree,
		dataType: reflect.TypeOf(""),
	}
}

// Insert adds an item to the shortener. If the value given is not of the same
// type of the elements in the original data set, Add panics.
func (s *Shortener) Insert(item interface{}) {
	v, ok := item.(reflect.Value)
	if !ok {
		v = reflect.ValueOf(item)
	}
	s.checkType(v)

	// Handle strings
	if v.Kind() == reflect.String {
		s.tree.Insert(v.String(), struct{}{})
		return
	}

	// Handle structs
	v = reflect.Indirect(v)
	s.tree.Insert(v.FieldByName(s.field).String(), struct{}{})
}

// Delete is used to remove an item from the data set.
func (s *Shortener) Delete(item interface{}) {
	v := reflect.ValueOf(item)
	s.checkType(v)

	// Handle strings
	if v.Kind() == reflect.String {
		s.tree.Delete(v.String())
		return
	}

	// Handle structs
	v = reflect.Indirect(v)
	s.tree.Delete(v.FieldByName(s.field).String())
}

// checkType asserts that the provided type matches what the shortener was
// created with, and panics if the types differ.
func (s *Shortener) checkType(v reflect.Value) {
	if v.Type() != s.dataType {
		panic(fmt.Sprintf("type must be %v, got %v", s.dataType, v.Type()))
	}
}

// min is the internal method used to retrieve the shortest possible string,
// given the length constraint.
func (s *Shortener) min(in string, l int) string {
	var result string
	for i := 0; ; i++ {
		// Add the next chunk of characters
		lidx := (i + 1) * l
		if lidx > len(in) {
			break
		}
		result += in[i*l : lidx]

		// Walk the tree. If we find more than a single result, then the
		// result would be ambiguous.
		var ambiguous, found bool
		s.tree.WalkPrefix(result, func(s string, _ interface{}) bool {
			if found {
				ambiguous = true
				return true
			}
			found = true
			return false
		})

		// If the prefix didn't match anything, then return early as the
		// prefix isn't in the data set.
		if !found {
			return ""
		}

		// If multiple entries were found for the prefix, continue to add more
		// characters to disambiguate.
		if ambiguous {
			continue
		}

		// We got an unambiguous result, so return it
		return result
	}
	return ""
}

// ShortestChunk is used to return the shortest substring in the chunk size
// provided. This means the minimum returned length is l, and the max is a
// multiple thereof. This is useful for keeping churn rate low with a
// frequently changing data set.
//
// If the result is an empty string, then shortening would create an ambiguous
// result (non-unique in the set).
func (s *Shortener) ShortestChunk(in string, l int) string {
	return s.min(in, l)
}

// Shortest is used to return the shortest possible unique match from the
// data set. If the result is an empty string, then shortening would create
// an ambiguous result (non-unique in the set).
func (s *Shortener) Shortest(in string) string {
	return s.min(in, 1)
}

// Expand is used to look up the full value of a given short string in the data
// set. If the result is an empty string, then expanding is not possible due to
// either the provided prefix missing in the data set, or multiple entries
// sharing the same prefix.
func (s *Shortener) Expand(in string) string {
	var ambiguous bool
	var full string

	// Walk the prefix of the given short string. If a single entry is found we
	// can return safely, but if we find more then the lookup cannot resolve.
	s.tree.WalkPrefix(in, func(s string, _ interface{}) bool {
		if full != "" {
			ambiguous = true
			return true
		}
		full = s
		return false
	})

	// Check if we found multiple entries by the same prefix.
	if ambiguous {
		return ""
	}

	// A single match was found, so return it.
	return full
}
