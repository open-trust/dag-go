# dag-go
A DAG, Directed acyclic graph implementation in golang.

[![License](https://img.shields.io/badge/license-MIT-blue)](https://raw.githubusercontent.com/open-trust/dag-go/master/LICENSE)

https://en.wikipedia.org/wiki/Directed_acyclic_graph

## Documentation
https://pkg.go.dev/github.com/open-trust/dag-go

![DAG](https://github.com/open-trust/dag-go/blob/main/dag-go.png?raw=true)

## Package

### Import
```go
import (
	daggo "github.com/open-trust/dag-go"
)
```

### Example
```go
package main

import (
	"encoding/json"
	"fmt"
	"time"

	otgo "github.com/open-trust/ot-go-lib"
)

type V string

func (v V) ID() string {
	return string(v)
}
func (v V) Attrs() daggo.Attrs {
	return daggo.Attrs{A(v)}
}

type A string

func (a A) ID() string {
	return string(a)
}

func main() {
  d := daggo.New()
  _ = d.AddEdge(V("a"), V("b"), 10)
  _ = d.AddEdge(V("a"), V("c"), 3)
  _ = d.AddEdge(V("a"), V("d"), 10)
  _ = d.AddEdge(V("a"), V("e"), 10)
  _ = d.AddEdge(V("b"), V("d"), 10)
  _ = d.AddEdge(V("c"), V("d"), 100)
  _ = d.AddEdge(V("c"), V("e"), 3)
  _ = d.AddEdge(V("d"), V("e"), 10)
  _ = d.AddEdge(V("x"), V("b"), 10)
  _ = d.AddEdge(V("d"), V("y"), 10)

  fmt.Println(d.Vertices()) // a, b, c, d, e, x, y
  fmt.Println(d.StartingVertices()) // x, a
  fmt.Println(d.EndingVertices()) // e, y
  fmt.Println(d.ToVertices(V("a"))) // b, c, e
  fmt.Println(d.FromVertices(V("e"))) // a, c, d
  fmt.Println(d.ReachDAG(V("a")))
  fmt.Println(d.CloseDAG(V("a"), V("e")))
  fmt.Println(d.ReduceDAG(V("a"), V("e")))
  fmt.Println(d.Reverse())
  fmt.Println(d.Shortest(V("a"), V("e"), false)) // a, e
  fmt.Println(d.Shortest(V("a"), V("e"), true)) // a, c, e
  fmt.Println(d.Longest(V("a"), V("e"), false)) // [a, c, d, e] or [a, b, d, e]
  fmt.Println(d.Longest(V("a"), V("e"), true)) // a, c, d, e

  fmt.Println(d.CloseDAG(V("a"), V("e")).
    Iterate(V("a"), nil, func(v daggo.Vertice, w int, acc daggo.Attrs) daggo.Attrs {
      return append(acc, v.Attrs()...)
    }))
}
```
