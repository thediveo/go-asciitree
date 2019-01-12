package asciitree_test

import (
	"fmt"

	asciitree "github.com/thediveo/go-asciitree"
)

func Example_sorted() {
	// user-defined tree data structure with asciitree-related field tags.
	type node struct {
		Label    string   `asciitree:"label"`
		Props    []string `asciitree:"properties"`
		Children []node   `asciitree:"children"`
	}
	// set up a tree of nodes.
	rootb := node{
		Label: "beta root",
		Children: []node{
			node{Label: "foo", Props: []string{"childish"}},
			node{Label: "alpha", Children: []node{
				node{Label: "grandchild 2"},
				node{Label: "grandchild 1", Props: []string{"very childish"}},
			}},
			node{Label: "bar"},
		},
	}
	roota := node{
		Label: "alpha root",
		Children: []node{
			node{Label: "alphachild"},
		},
	}
	// create a new visitor and tell it to sort the nodes by label, and also
	// to sort the properties.
	sortingVisitor := asciitree.NewMapStructVisitor(true, true)
	// render the tree(s) into a string and print it.
	fmt.Println(asciitree.Render([]node{rootb, roota}, sortingVisitor, asciitree.LineTreeStyler))
	// Output:
	// alpha root
	// └─ alphachild
	// beta root
	// ├─ alpha
	// │  ├─ grandchild 1
	// │  │     • very childish
	// │  └─ grandchild 2
	// ├─ bar
	// └─ foo
	//       • childish
}
