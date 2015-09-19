package shortstr

import (
	"testing"
)

func checkBasic(t *testing.T, s *Shortener) {
	if res := s.Shortest("aaaa"); res != "aaa" {
		t.Fatalf("expect %q, got %q", "aaa", res)
	}
	if res := s.Shortest("nope"); res != "" {
		t.Fatalf("expect %q, got %q", "", res)
	}
	if res := s.ShortestChunk("aaaa", 4); res != "aaaa" {
		t.Fatalf("expect %q, got %q", "aaaa", res)
	}
	if res := s.ShortestChunk("nope", 4); res != "" {
		t.Fatalf("expect %q, got %q", "", res)
	}
	if res := s.Expand("aaa"); res != "aaaa" {
		t.Fatalf("expect %q, got %q", "aaaa", res)
	}
	if res := s.Expand("nope"); res != "" {
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

func TestStructPointers(t *testing.T) {
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
	New([]thing{}, "")
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
