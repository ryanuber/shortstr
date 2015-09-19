shortstr
========

`shortstr` is a Golang library which can be used to shorten strings or struct
field values. It provides a two-way mapping which can take a long string and
return a shorter version, as well as taking a shortened string and returning
the full version.

Internally `shortstr` uses [radix trees]
(https://en.wikipedia.org/wiki/Radix_tree) to optimize prefix lookups. This
provides O(k) operations, which are very similar in overhead to a hash table.

Example usage
=============

Using custom structs:

```go
type Thing struct {
    Name string
}

func main() {
    data := []*Thing{
        &Thing{"aaaaaa"},
        &Thing{"aabbbb"},
        &Thing{"aacccc"},
    }
    s := shortstr.New(data, "Name")
    println(s.Shortest("aaaaaa"))         // "aaa"
    println(s.ShortestChunk("aabbbb", 4)) // "aabb"
    println(s.Expand("aac"))              // "aacccc"
}
```

Using string slices:

```go
func main() {
    data := []string{
        "aaaaaa",
        "aabbbb",
        "aacccc",
    }
    s := shortstr.NewStrings(data)
    println(s.Shortest("aaaaaa"))         // "aaa"
    println(s.ShortestChunk("aabbbb", 4)) // "aabb"
    println(s.Expand("aac"))              // "aacccc"
}
```
