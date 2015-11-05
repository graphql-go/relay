package relay

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const PREFIX = "arrayconnection:"

/*
A simple function that accepts an array and connection arguments, and returns
a connection object for use in GraphQL. It uses array offsets as pagination,
so pagination will only work if the array is static.
*/

func ConnectionFromArray(data []interface{}, args ConnectionArguments) *Connection {
	edges := []*Edge{}
	for index, value := range data {
		edges = append(edges, &Edge{
			Cursor: offsetToCursor(index),
			Node:   value,
		})
	}

	// slice with cursors
	afterOffset := getOffset(args.After, -1)
	beforeOffset := getOffset(args.Before, len(edges)+1)

	begin := int(math.Max(float64(afterOffset), -1) + 1)
	end := int(math.Min(float64(beforeOffset), float64(len(edges))))
	if begin > end {
		return NewConnection()
	}

	edges = edges[begin:end]
	if len(edges) == 0 {
		return NewConnection()
	}

	// save the pre-slice cursors
	firstPresliceCursor := edges[0].Cursor
	lastPresliceCursor := edges[len(edges)-1:][0].Cursor

	// slice with limits
	if args.First >= 0 {
		first := int(math.Min(float64(args.First), float64(len(edges))))
		edges = edges[0:first]
	}

	if args.Last >= 0 {
		last := int(math.Min(float64(args.Last), float64(len(edges))))
		edges = edges[len(edges)-last:]
	}

	if len(edges) == 0 {
		return NewConnection()
	}
	firstEdge := edges[0]
	lastEdge := edges[len(edges)-1:][0]
	hasPreviousPage := false
	if firstEdge.Cursor != firstPresliceCursor {
		hasPreviousPage = true
	}
	hasNextPage := false
	if lastEdge.Cursor != lastPresliceCursor {
		hasNextPage = true
	}

	conn := NewConnection()
	conn.Edges = edges
	conn.PageInfo = PageInfo{
		StartCursor:     firstEdge.Cursor,
		EndCursor:       lastEdge.Cursor,
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
	}
	return conn
}

// Creates the cursor string from an offset
func offsetToCursor(offset int) ConnectionCursor {
	str := fmt.Sprintf("%v%v", PREFIX, offset)
	return ConnectionCursor(base64.StdEncoding.EncodeToString([]byte(str)))
}

// Re-derives the offset from the cursor string.
func cursorToOffset(cursor ConnectionCursor) (int, error) {

	str := ""
	b, err := base64.StdEncoding.DecodeString(string(cursor))
	if err == nil {
		str = string(b)
	}
	str = strings.Replace(str, PREFIX, "", -1)
	offset, err := strconv.Atoi(str)
	if err != nil {
		return 0, errors.New("Invalid cursor")
	}
	return offset, nil
}

// Return the cursor associated with an object in an array.
func CursorForObjectInConnection(data []interface{}, object interface{}) ConnectionCursor {
	offset := -1
	for i, d := range data {
		// TODO: better object comparison
		if d == object {
			offset = i
			break
		}
	}
	if offset == -1 {
		return ""
	}
	return offsetToCursor(offset)
}

func getOffset(cursor ConnectionCursor, defaultOffset int) int {
	if cursor == "" {
		return defaultOffset
	}
	offset, err := cursorToOffset(cursor)
	if err != nil {
		return defaultOffset
	}
	return offset
}
