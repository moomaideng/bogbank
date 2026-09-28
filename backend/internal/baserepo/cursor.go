package baserepo

// TODO: keyset pagination is deferred until a list endpoint needs it.

// CursorInput is a Relay-style page request. First and After walk forward.
// Last and Before walk backward. Cursors are opaque strings the repository
// defines.
type CursorInput struct {
	First  *int
	Last   *int
	Before *string
	After  *string
}

// CursorPageInfo describes the page around a cursor result.
type CursorPageInfo struct {
	HasNextPage     bool
	HasPreviousPage bool
	StartCursor     string
	EndCursor       string
}

// CursorEdge is one node plus the cursor that addresses it.
type CursorEdge[T any] struct {
	Node   T
	Cursor string
}

// CursorPage is one cursor result.
type CursorPage[T any] struct {
	Edges    []CursorEdge[T]
	PageInfo CursorPageInfo
}
