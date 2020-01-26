package asciitree_test

import (
	"fmt"

	asciitree "github.com/thediveo/go-asciitree"
)

func ExampleRender() {
	// user-defined tree data structure with asciitree-related field tags.
	type tree struct {
		Label    string `asciitree:"label"`
		Children []tree `asciitree:"children"`
	}
	// set up a tree of nodes.
	root := tree{
		Label: "root",
		Children: []tree{
			{Label: "child 1"},
			{Label: "child 2", Children: []tree{
				{Label: "grandchild 1"},
				{Label: "grandchild 2"},
			}},
			{Label: "child 3"},
		},
	}
	// render the tree into a string and print it.
	fmt.Println(
		asciitree.Render(root,
			asciitree.DefaultVisitor,
			asciitree.DefaultTreeStyler))
	// Output:
	// root
	// +- child 1
	// +- child 2
	// |  +- grandchild 1
	// |  `- grandchild 2
	// `- child 3
}

func ExampleRenderPlain() {
	// user-defined tree data structure with asciitree-related field tags.
	type tree struct {
		Label    string `asciitree:"label"`
		Children []tree `asciitree:"children"`
	}
	// set up a tree of nodes.
	root := tree{
		Label: "root",
		Children: []tree{
			{Label: "child 1"},
			{Label: "child 2", Children: []tree{
				{Label: "grandchild 1"},
				{Label: "grandchild 2"},
			}},
			{Label: "child 3"},
		},
	}
	// render the tree into a string and print it.
	fmt.Println(asciitree.RenderPlain(root))
	// Output:
	// root
	// +- child 1
	// +- child 2
	// |  +- grandchild 1
	// |  `- grandchild 2
	// `- child 3
}

func ExampleRenderFancy() {
	// user-defined tree data structure with asciitree-related field tags.
	type tree struct {
		Label    string `asciitree:"label"`
		Children []tree `asciitree:"children"`
	}
	// set up a tree of nodes.
	root := tree{
		Label: "root",
		Children: []tree{
			{Label: "child 1"},
			{Label: "child 2", Children: []tree{
				{Label: "grandchild 1"},
				{Label: "grandchild 2"},
			}},
			{Label: "child 3"},
		},
	}
	// render the tree into a string and print it.
	fmt.Println(asciitree.RenderFancy(root))
	// Output:
	// root
	// ├─ child 1
	// ├─ child 2
	// │  ├─ grandchild 1
	// │  └─ grandchild 2
	// └─ child 3
}
