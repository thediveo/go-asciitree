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
			{Label: "foo", Props: []string{"childish"}},
			{Label: "alpha", Children: []node{
				{Label: "grandchild 2"},
				{Label: "grandchild 1", Props: []string{"very childish"}},
			}},
			{Label: "bar"},
		},
	}
	roota := node{
		Label: "alpha root",
		Children: []node{
			{Label: "alphachild"},
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
