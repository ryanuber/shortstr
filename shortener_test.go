package shortstr

import (
	"testing"
)

func checkBasic(t *testing.T, s *Shortener) {
	// Returns the shortest possible string
	if res := s.Shortest("aaaa"); res != "aaa" {
		t.Fatalf("expect %q, got %q", "aaa", res)
	}

	// Returns nothing on no match
	if res := s.Shortest("nope"); res != "" {
		t.Fatalf("expect %q, got %q", "", res)
	}

	// Returns nothing if ambiguous
	if res := s.Shortest("aa"); res != "" {
		t.Fatalf("expect %q, got %q", "", res)
	}

	// Returns the shortest chunked match
	if res := s.ShortestChunk("aaaa", 2); res != "aaaa" {
		t.Fatalf("expect %q, got %q", "aaaa", res)
	}

	// Returns nothing on no match
	if res := s.ShortestChunk("nope", 2); res != "" {
		t.Fatalf("expect %q, got %q", "", res)
	}

	// Returns nothing if ambiguous
	if res := s.ShortestChunk("aa", 2); res != "" {
		t.Fatalf("expect %q, got %q", "", res)
	}

	// Returns the expanded version of the string
	if res := s.Expand("aaa"); res != "aaaa" {
		t.Fatalf("expect %q, got %q", "aaaa", res)
	}

	// Returns nothing if no matches
	if res := s.Expand("nope"); res != "" {
		t.Fatalf("expect %q, got %q", "", res)
	}

	// Returns nothing if multiple matches
	if res := s.Expand("aa"); res != "" {
		t.Fatalf("expect %q, got %q", "", res)
	}
}

func TestStrings(t *testing.T) {
	s := NewStrings([]string{"aaaa", "aabb", "aacc"})
	checkBasic(t, s)
}

func TestStructValues(t *testing.T) {
	type thing struct {
		value string
	}
	s := New([]thing{
		thing{"aaaa"},
		thing{"aabb"},
		thing{"aacc"},
	}, "value")
	checkBasic(t, s)
}

func TestStructPtrs(t *testing.T) {
	type thing struct {
		value string
	}
	s := New([]*thing{
		&thing{"aaaa"},
		&thing{"aabb"},
		&thing{"aacc"},
	}, "value")
	checkBasic(t, s)
}

func TestInvalidType(t *testing.T) {
	defer func(t *testing.T) {
		r := recover()
		if r == nil || r.(string) != "not a slice" {
			t.Fatalf("expected slice error, got: %#v", r)
		}
	}(t)
	New("nope", "nope")
}

func TestInvalidSliceType(t *testing.T) {
	defer func(t *testing.T) {
		r := recover()
		if r == nil || r.(string) != "not a struct slice" {
			t.Fatalf("expected slice error, got: %#v", r)
		}
	}(t)
	New([]string{}, "nope")
}

func TestInvalidSliceTypePtr(t *testing.T) {
	defer func(t *testing.T) {
		r := recover()
		if r == nil || r.(string) != "not a struct slice" {
			t.Fatalf("expected slice error, got: %#v", r)
		}
	}(t)
	New([]*string{}, "nope")
}

func TestInvalidFieldName(t *testing.T) {
	defer func(t *testing.T) {
		r := recover()
		if r == nil || r.(string) != "invalid struct field" {
			t.Fatalf("expected field name error, got: %#v", r)
		}
	}(t)
	type thing struct{}
	New([]thing{}, "nope")
}

func TestInvalidFieldNamePtr(t *testing.T) {
	defer func(t *testing.T) {
		r := recover()
		if r == nil || r.(string) != "invalid struct field" {
			t.Fatalf("expected field name error, got: %#v", r)
		}
	}(t)
	type thing struct{}
	New([]*thing{}, "nope")
}

func TestInvalidFieldType(t *testing.T) {
	defer func(t *testing.T) {
		r := recover()
		if r == nil || r.(string) != "struct field must be type string" {
			t.Fatalf("expected field type error, got: %#v", r)
		}
	}(t)
	type thing struct {
		bad int
	}
	New([]thing{}, "bad")
}

func TestInvalidFieldTypePtr(t *testing.T) {
	defer func(t *testing.T) {
		r := recover()
		if r == nil || r.(string) != "struct field must be type string" {
			t.Fatalf("expected field type error, got: %#v", r)
		}
	}(t)
	type thing struct {
		bad int
	}
	New([]*thing{}, "bad")
}

func TestAddInvalidType(t *testing.T) {
	defer func(t *testing.T) {
		r := recover()
		if r == nil || r.(string) != "type must be string, got struct {}" {
			t.Fatalf("expected data type error, got: %#v", r)
		}
	}(t)
	s := NewStrings([]string{})
	s.Add(struct{}{})
}

func TestAdd(t *testing.T) {
	s := NewStrings([]string{"foo", "bar"})
	s.Add("baz")
	if _, ok := s.tree.Get("baz"); !ok {
		t.Fatalf("missing added entry")
	}
}

func TestRemove(t *testing.T) {
	s := NewStrings([]string{"foo", "bar"})
	s.Delete("bar")
	if _, ok := s.tree.Get("bar"); ok {
		t.Fatalf("should remove entry")
	}
}
