package regtree

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type storeTreeEntry struct {
	key    string
	data   string
	params int
}

func TestTreeAdd(t *testing.T) {
	tests := []struct {
		id      string
		entries []storeTreeEntry
	}{
		{
			"all static",
			[]storeTreeEntry{
				{"/gopher/bumper.png", "1", 0},
				{"/gopher/bumper192x108.png", "2", 0},
				{"/gopher/doc.png", "3", 0},
				{"/gopher/bumper320x180.png", "4", 0},
				{"/gopher/docpage.png", "5", 0},
				{"/gopher/doc", "7", 0},
			},
		},
		{
			"parametric",
			[]storeTreeEntry{
				{"/users/:id", "11", 1},
				{"/users/:id/profile", "12", 1},
				{"/users/:id/:accnt(\\d+)/address", "13", 2},
				{"/users/:id/age", "14", 1},
				{"/users/:id/:accnt(\\d+)", "15", 2},
			},
		},
		{
			"corner cases",
			[]storeTreeEntry{
				{"/users/:id/test/:name", "101", 2},
				{"/users/abc/:id/:name", "102", 2},
			},
		},
	}

	Convey("tree add", t, func() {
		for _, test := range tests {
			h := NewTree("/", nil)
			Convey(test.id, func() {
				for _, entry := range test.entries {
					node := h.Add(entry.key, entry.data)
					if node == nil {
						fmt.Printf("nil node: %#v", entry)
					}
					So(node, ShouldNotBeNil)
					if len(node.params) != entry.params {
						fmt.Printf("error node: %#v\n", node)
					}
					So(len(node.params), ShouldEqual, entry.params)
				}
			})
		}
	})
}

func TestStoreGet(t *testing.T) {
	pairs := []struct {
		key, value string
	}{
		{"/gopher/bumper.png", "1"},
		{"/gopher/bumper192x108.png", "2"},
		{"/gopher/doc.png", "3"},
		{"/gopher/bumper320x180.png", "4"},
		{"/gopher/docpage.png", "5"},
		{"/gopher/doc", "7"},
		{"/users/:id", "8"},
		{"/users/:id/profile", "9"},
		{"/users/:id/:accnt(\\d+)/address", "10"},
		{"/users/:id/age", "11"},
		{"/users/:id/:accnt(\\d+)", "12"},
		{"/users/:id/test/:name", "13"},
		{"/users/abc/:id/:name", "14"},
		{"/all/*", "16"},
	}
	h := NewTree("/", nil)
	for _, pair := range pairs {
		h.Add(pair.key, pair.value)
	}

	tests := []struct {
		key    string
		value  interface{}
		params string
	}{
		{"/gopher/bumper.png", "1", ""},
		{"/gopher/bumper192x108.png", "2", ""},
		{"/gopher/doc.png", "3", ""},
		{"/gopher/bumper320x180.png", "4", ""},
		{"/gopher/docpage.png", "5", ""},
		{"/gopher/doc", "7", ""},
		{"/users/abc", "8", "id:abc,"},
		{"/users/abc/profile", "9", "id:abc,"},
		{"/users/abc/123/address", "14", "id:123,name:address,"},
		{"/users/abcd/age", "11", "id:abcd,"},
		{"/users/abc/123", "12", "id:abc,accnt:123,"},
		{"/users/abc/test/123", "14", "id:test,name:123,"},
		{"/users/abc/xyz/123", "14", "id:xyz,name:123,"},
		{"/g", nil, ""},
		{"/all", nil, ""},
		{"/all/", "16", ":,"},
		{"/all/abc", "16", ":abc,"},
		{"/users/abc/xyz", nil, ""},
	}
	Convey("tree get", t, func() {
		for _, test := range tests {
			node, values := h.Get(test.key)
			if test.value == nil {
				So(node, ShouldBeNil)
			}
			if node == nil {
				t.Logf("found nil: %#v\n", test)
				continue
			}
			if node.val != test.value {
				fmt.Printf("node: %#v\n", node)
			}
			So(node.val, ShouldEqual, test.value)
			params := ""
			if len(node.params) > 0 {
				for i, name := range node.params {
					params += fmt.Sprintf("%v:%v,", name, values[i])
				}
			}
			So(params, ShouldEqual, test.params)
		}
	})
}
