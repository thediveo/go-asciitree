package asciitree_test

import (
	"fmt"

	asciitree "github.com/thediveo/go-asciitree"
)

func Example() {
	// user-defined tree data structure with asciitree-related field tags.
	type node struct {
		Label    string   `asciitree:"label"`
		Props    []string `asciitree:"properties"`
		Children []node   `asciitree:"children"`
	}
	// set up a tree of nodes.
	root := node{
		Label: "root",
		Children: []node{
			node{Label: "child 1", Props: []string{"childish"}},
			node{Label: "child 2", Children: []node{
				node{Label: "grandchild 1", Props: []string{"very childish"}},
				node{Label: "grandchild 2"},
			}},
			node{Label: "child 3"},
		},
	}
	// render the tree into a string and print it.
	fmt.Println(asciitree.RenderFancy(root))
	// Output:
	// root
	// ├─ child 1
	// │     • childish
	// ├─ child 2
	// │  ├─ grandchild 1
	// │  │     • very childish
	// │  └─ grandchild 2
	// └─ child 3
}
