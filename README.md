shortstr
========

[![GoDoc](https://godoc.org/github.com/ryanuber/shortstr?status.svg)](https://godoc.org/github.com/ryanuber/shortstr)

`shortstr` is a Golang library which can be used to shorten strings or struct
field values while still providing reasonable uniqueness guarantees. It
provides a two-way mapping which can take a long string and return a shorter
version, as well as taking a shortened string and returning the full version.

A practical example would be shortening human-facing UUID's. A UUID is a very
convenient way to assign a unique ID to an object, but because of its uniqueness
guarantees, it suffers from human grok-ability, and is near impossible to
remember when presented with long lists. Often times computer operators will
memorize a short segment, perhaps the first 4 to 8 characters, and that is
usually enough to uniquely identify the object in question. This library aims to
make using these shortened forms easy while still retaining uniqueness
guarantees and the ability to map abbreviated values back to their original,
full-sized values.

Performance
===========

Internally `shortstr` uses [radix trees]
(https://en.wikipedia.org/wiki/Radix_tree) to optimize prefix lookups. A radix
tree provides O(k) operations, which have search times comparable to a hash
table. `shortstr` may perform multiple lookups to determine the shortest
possible match for a string. The implementation used in this library is
[go-radix](https://github.com/armon/go-radix).

Example usage
=============

```go
type Thing struct {
    Name string
}

data := []Thing{
    Thing{"aaaaaaaa"},
    Thing{"aaaabbbb"},
}

// Create the shortener with an initial data set
s := shortstr.New(data, "Name")

// Basic shorten/expand
println(s.Shortest("aaaabbbb"))         // "aaaab"
println(s.ShortestChunk("aaaabbbb", 6)) // "aaaabb"
println(s.Expand("aaaab"))              // "aaaabbbb"

// Add more data items to the set
s.Add(Thing{"aaaabbcc"})

// Remove data items
s.Remove(Thing{"aaaabbcc"})
```

It is also possible to create the shortener using a slice of strings:

```go
data := []string{
    "aaaaaa",
    "aabbbb",
    "aacccc",
}
s := shortstr.NewStrings(data)
```
