package daggo_test

import (
	"testing"

	daggo "github.com/open-trust/dag-go"
	"github.com/stretchr/testify/assert"
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

func TestDAG(t *testing.T) {
	t.Run("should work", func(t *testing.T) {
		assert := assert.New(t)

		d := daggo.New()
		assert.Equal(0, d.Len())
		assert.True(d.GetVertice("a") == nil)
		assert.Equal(daggo.Vertices{}, d.Vertices())
		assert.Equal(daggo.Vertices{}, d.StartingVertices())
		assert.Equal(daggo.Vertices{}, d.EndingVertices())
		assert.Equal(daggo.Vertices{}, d.ToVertices(nil))
		assert.Equal(daggo.Vertices{}, d.FromVertices(nil))
		assert.Equal(daggo.IDs{}, d.VerticeIDs())
		assert.Equal(daggo.Attrs{}, d.Attrs())

		assert.NotNil(d.AddEdge(nil, nil, 0))
		assert.NotNil(d.AddEdge(V("a"), nil, 0))
		assert.NotNil(d.AddEdge(V("a"), V("a"), 0))
		assert.Nil(d.AddEdge(V("a"), V("b"), 0))

		assert.Equal(2, d.Len())
		assert.Equal("a", d.GetVertice("a").ID())
		assert.Equal("b", d.GetVertice("b").ID())
		assert.Equal(nil, d.GetVertice("c"))
		assert.Equal(daggo.Vertices{V("a")}, d.StartingVertices())
		assert.Equal(daggo.Vertices{V("b")}, d.EndingVertices())
		assert.Equal(daggo.Vertices{V("b")}, d.ToVertices(V("a")))
		assert.Equal(daggo.Vertices{V("a")}, d.FromVertices(V("b")))
		assert.Equal(daggo.IDs{"a", "b"}, d.VerticeIDs().Sort())

		assert.Nil(d.AddEdge(V("a"), V("c"), 0))
		assert.Nil(d.AddEdge(V("x"), V("b"), 0))
		assert.Equal(daggo.Vertices{V("a"), V("x")}, d.StartingVertices().Sort())
		assert.Equal(daggo.Vertices{V("b"), V("c")}, d.EndingVertices().Sort())
		assert.Equal(daggo.Vertices{V("b"), V("c")}, d.ToVertices(V("a")).Sort())
		assert.Equal(daggo.Vertices{V("a"), V("x")}, d.FromVertices(V("b")).Sort())

		assert.Nil(d.AddEdge(V("a"), V("x"), 0))
		assert.NotNil(d.AddEdge(V("b"), V("a"), 0))
		assert.Equal(daggo.Vertices{V("a")}, d.StartingVertices().Sort())
		assert.Equal(daggo.Vertices{V("b"), V("c")}, d.EndingVertices().Sort())

		d.RemoveEdge(V("a"), V("b"))
		assert.NotNil(d.AddEdge(V("b"), V("a"), 0))
		assert.Equal(daggo.Vertices{V("a")}, d.StartingVertices().Sort())
		assert.Equal(daggo.Vertices{V("b"), V("c")}, d.EndingVertices().Sort())
		assert.Equal(daggo.Vertices{V("c"), V("x")}, d.ToVertices(V("a")).Sort())
		assert.Equal(daggo.Vertices{V("x")}, d.FromVertices(V("b")).Sort())
	})

	t.Run("DAG.ReachDAG", func(t *testing.T) {
		assert := assert.New(t)

		d := daggo.New()
		assert.Nil(d.AddEdge(V("a"), V("b"), 0))
		assert.Nil(d.AddEdge(V("a"), V("c"), 0))
		assert.Nil(d.AddEdge(V("a"), V("d"), 0))
		assert.Nil(d.AddEdge(V("a"), V("e"), 0))
		assert.Nil(d.AddEdge(V("b"), V("d"), 0))
		assert.Nil(d.AddEdge(V("c"), V("d"), 0))
		assert.Nil(d.AddEdge(V("c"), V("e"), 0))
		assert.Nil(d.AddEdge(V("d"), V("e"), 0))
		assert.Nil(d.AddEdge(V("x"), V("b"), 0))
		assert.Nil(d.AddEdge(V("d"), V("y"), 0))

		x := daggo.New()
		assert.Nil(x.AddEdge(V("a"), V("b"), 0))
		assert.Nil(x.AddEdge(V("a"), V("c"), 0))
		assert.Nil(x.AddEdge(V("a"), V("d"), 0))
		assert.Nil(x.AddEdge(V("a"), V("e"), 0))
		assert.Nil(x.AddEdge(V("b"), V("d"), 0))
		assert.Nil(x.AddEdge(V("c"), V("d"), 0))
		assert.Nil(x.AddEdge(V("c"), V("e"), 0))
		assert.Nil(x.AddEdge(V("d"), V("e"), 0))
		assert.Nil(x.AddEdge(V("d"), V("y"), 0))
		assert.True(x.Equal(d.ReachDAG(V("a"))))

		assert.Nil(x.AddEdge(V("a"), V("b"), 1))
		assert.False(x.Equal(d.ReachDAG(V("a"))))
	})

	t.Run("DAG.CloseDAG", func(t *testing.T) {
		assert := assert.New(t)

		d := daggo.New()
		assert.Nil(d.AddEdge(V("a"), V("b"), 0))
		assert.Nil(d.AddEdge(V("a"), V("c"), 0))
		assert.Nil(d.AddEdge(V("a"), V("d"), 0))
		assert.Nil(d.AddEdge(V("a"), V("e"), 0))
		assert.Nil(d.AddEdge(V("b"), V("d"), 0))
		assert.Nil(d.AddEdge(V("c"), V("d"), 0))
		assert.Nil(d.AddEdge(V("c"), V("e"), 0))
		assert.Nil(d.AddEdge(V("d"), V("e"), 0))
		assert.Nil(d.AddEdge(V("x"), V("b"), 0))
		assert.Nil(d.AddEdge(V("d"), V("y"), 0))

		x := daggo.New()
		assert.Nil(x.AddEdge(V("a"), V("b"), 0))
		assert.Nil(x.AddEdge(V("a"), V("c"), 0))
		assert.Nil(x.AddEdge(V("a"), V("d"), 0))
		assert.Nil(x.AddEdge(V("a"), V("e"), 0))
		assert.Nil(x.AddEdge(V("b"), V("d"), 0))
		assert.Nil(x.AddEdge(V("c"), V("d"), 0))
		assert.Nil(x.AddEdge(V("c"), V("e"), 0))
		assert.Nil(x.AddEdge(V("d"), V("e"), 0))
		assert.True(x.Equal(d.CloseDAG(V("a"), V("e"))))

		assert.Nil(d.AddEdge(V("z"), V("a"), 0))
		assert.True(x.Equal(d.CloseDAG(V("a"), V("e"))))

		x = daggo.New()
		assert.Nil(x.AddEdge(V("a"), V("b"), 0))
		assert.Nil(x.AddEdge(V("a"), V("c"), 0))
		assert.Nil(x.AddEdge(V("a"), V("d"), 0))
		assert.Nil(x.AddEdge(V("b"), V("d"), 0))
		assert.Nil(x.AddEdge(V("c"), V("d"), 0))
		assert.True(x.Equal(d.CloseDAG(V("a"), V("d"))))
	})

	t.Run("DAG.ReduceDAG", func(t *testing.T) {
		assert := assert.New(t)

		d := daggo.New()
		assert.Nil(d.AddEdge(V("a"), V("b"), 0))
		assert.Nil(d.AddEdge(V("a"), V("c"), 0))
		assert.Nil(d.AddEdge(V("a"), V("d"), 0))
		assert.Nil(d.AddEdge(V("a"), V("e"), 0))
		assert.Nil(d.AddEdge(V("b"), V("d"), 0))
		assert.Nil(d.AddEdge(V("c"), V("d"), 0))
		assert.Nil(d.AddEdge(V("c"), V("e"), 0))
		assert.Nil(d.AddEdge(V("d"), V("e"), 0))
		assert.Nil(d.AddEdge(V("x"), V("b"), 0))
		assert.Nil(d.AddEdge(V("d"), V("y"), 0))

		x := daggo.New()
		assert.Nil(x.AddEdge(V("a"), V("b"), 0))
		assert.Nil(x.AddEdge(V("a"), V("c"), 0))
		assert.Nil(x.AddEdge(V("b"), V("d"), 0))
		assert.Nil(x.AddEdge(V("c"), V("d"), 0))
		assert.Nil(x.AddEdge(V("d"), V("e"), 0))
		assert.True(x.Equal(d.ReduceDAG(V("a"), V("e"))))

		assert.Nil(d.AddEdge(V("z"), V("a"), 0))
		assert.True(x.Equal(d.ReduceDAG(V("a"), V("e"))))

		x = daggo.New()
		assert.Nil(x.AddEdge(V("a"), V("b"), 0))
		assert.Nil(x.AddEdge(V("a"), V("c"), 0))
		assert.Nil(x.AddEdge(V("b"), V("d"), 0))
		assert.Nil(x.AddEdge(V("c"), V("d"), 0))
		assert.True(x.Equal(d.ReduceDAG(V("a"), V("d"))))
	})

	t.Run("DAG.Reverse", func(t *testing.T) {
		assert := assert.New(t)

		d := daggo.New()
		assert.Nil(d.AddEdge(V("a"), V("b"), 0))
		assert.Nil(d.AddEdge(V("a"), V("c"), 0))
		assert.Nil(d.AddEdge(V("a"), V("d"), 0))
		assert.Nil(d.AddEdge(V("a"), V("e"), 0))
		assert.Nil(d.AddEdge(V("b"), V("d"), 0))
		assert.Nil(d.AddEdge(V("c"), V("d"), 0))
		assert.Nil(d.AddEdge(V("c"), V("e"), 0))
		assert.Nil(d.AddEdge(V("d"), V("e"), 0))
		assert.Nil(d.AddEdge(V("x"), V("b"), 0))
		assert.Nil(d.AddEdge(V("d"), V("y"), 0))

		x := daggo.New()
		assert.Nil(x.AddEdge(V("b"), V("a"), 0))
		assert.Nil(x.AddEdge(V("c"), V("a"), 0))
		assert.Nil(x.AddEdge(V("d"), V("a"), 0))
		assert.Nil(x.AddEdge(V("e"), V("a"), 0))
		assert.Nil(x.AddEdge(V("d"), V("b"), 0))
		assert.Nil(x.AddEdge(V("d"), V("c"), 0))
		assert.Nil(x.AddEdge(V("e"), V("c"), 0))
		assert.Nil(x.AddEdge(V("e"), V("d"), 0))
		assert.Nil(x.AddEdge(V("b"), V("x"), 0))
		assert.Nil(x.AddEdge(V("y"), V("d"), 0))

		assert.True(x.Equal(d.Reverse()))
	})

	t.Run("DAG.Shortest & Longest", func(t *testing.T) {
		assert := assert.New(t)

		d := daggo.New()
		assert.Nil(d.AddEdge(V("a"), V("x"), 10))
		assert.Nil(d.AddEdge(V("x"), V("b"), 10))
		assert.Nil(d.AddEdge(V("a"), V("c"), 10))
		assert.Nil(d.AddEdge(V("a"), V("d"), 10))
		assert.Nil(d.AddEdge(V("a"), V("e"), 10))
		assert.Nil(d.AddEdge(V("b"), V("d"), 10))
		assert.Nil(d.AddEdge(V("c"), V("d"), 10))
		assert.Nil(d.AddEdge(V("c"), V("e"), 10))
		assert.Nil(d.AddEdge(V("d"), V("e"), 10))
		assert.Nil(d.AddEdge(V("d"), V("y"), 10))

		assert.Equal(daggo.Vertices{V("a"), V("e")}, d.Shortest(V("a"), V("e"), false))
		assert.Equal(daggo.Vertices{V("a"), V("x"), V("b"), V("d"), V("e")}, d.Longest(V("a"), V("e"), false))

		assert.Nil(d.AddEdge(V("a"), V("c"), 3))
		assert.Nil(d.AddEdge(V("c"), V("e"), 3))
		assert.Equal(daggo.Vertices{V("a"), V("c"), V("e")}, d.Shortest(V("a"), V("e"), true))

		assert.Nil(d.AddEdge(V("c"), V("d"), 100))
		assert.Equal(daggo.Vertices{V("a"), V("c"), V("d"), V("e")}, d.Longest(V("a"), V("e"), true))
	})

	t.Run("DAG.Iterate", func(t *testing.T) {
		assert := assert.New(t)

		d := daggo.New()
		assert.Nil(d.AddEdge(V("a"), V("b"), 1))
		assert.Nil(d.AddEdge(V("a"), V("c"), 1))
		assert.Nil(d.AddEdge(V("a"), V("d"), 1))
		assert.Nil(d.AddEdge(V("a"), V("e"), 1))
		assert.Nil(d.AddEdge(V("b"), V("d"), 1))
		assert.Nil(d.AddEdge(V("c"), V("d"), 1))
		assert.Nil(d.AddEdge(V("c"), V("e"), 1))
		assert.Nil(d.AddEdge(V("d"), V("e"), 1))
		assert.Nil(d.AddEdge(V("x"), V("b"), 1))
		assert.Nil(d.AddEdge(V("d"), V("y"), 1))

		ws := 0
		attrs := d.CloseDAG(V("a"), V("e")).
			Iterate(V("a"), nil, func(v daggo.Vertice, w int, acc daggo.Attrs) daggo.Attrs {
				ws += w
				return append(acc, v.Attrs()...)
			})

		assert.Equal(10, ws)
		assert.Equal(daggo.Attrs{A("a"), A("a"), A("a"), A("a"), A("a"), A("b"), A("c"), A("c"), A("d"), A("d"), A("d"), A("e"), A("e"), A("e"), A("e"), A("e")}, attrs.Sort())
	})
}
