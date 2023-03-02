// Copyright 2018 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package asciitree

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type Node struct {
	Name       string   `asciitree:"label"`
	Properties []string `asciitree:"properties"`
	Subnodes   []*Node  `asciitree:"children"`
}

var _ = Describe("asciitree", func() {

	sortingVisitor := NewMapStructVisitor(true, true)

	rootnode1 := Node{
		Name:       "root1",
		Properties: []string{"foo", "bar"},
		Subnodes: []*Node{
			{Name: "1"},
			{
				Name: "2",
				Subnodes: []*Node{
					{Name: "2.1", Properties: []string{"whoooosh"}},
					{Name: "2.2"},
				},
			},
			{
				Name: "3",
				Subnodes: []*Node{
					{Name: "3.1"},
				},
			},
		},
	}
	rootnode2 := Node{
		Name: "root2",
		Subnodes: []*Node{
			{Name: "X"},
		},
	}
	type M map[string]interface{}
	rootmap := M{
		"roots": []M{
			{"label": "root1", "properties": []string{"foo", "bar"}, "children": []M{
				{"label": "1"},
				{"label": "2", "children": []M{
					{"label": "2.2"},
					{"label": "2.1", "properties": []string{"whoooosh"}},
				}},
				{"label": "3", "children": []M{
					{"label": "3.1"},
				}},
			}},
			{"label": "alpharot", "properties": []string{"z", "a"}},
		},
	}
	rootmap2 := M{"label": "root", "properties": []string{"pr"}, "children": []M{
		{"label": "1", "properties": []string{"p1", "p2"}},
	}}

	ts := NewTreeStyler(LineStyle)
	ts.ChildIndent = 4
	ts.PropIndent = 3

	It("renders root slices of nodes", func() {
		text := Render([]Node{rootnode1, rootnode2}, DefaultVisitor, ts)
		Expect(text).To(Equal(`root1
│  • foo
│  • bar
├── 1
├── 2
│   ├── 2.1
│   │      • whoooosh
│   └── 2.2
└── 3
    └── 3.1
root2
└── X
`))
	})

	It("renders sorted slices of nodes", func() {
		text := Render([]Node{rootnode2, rootnode1}, sortingVisitor, ts)
		Expect(text).To(Equal(`root1
│  • bar
│  • foo
├── 1
├── 2
│   ├── 2.1
│   │      • whoooosh
│   └── 2.2
└── 3
    └── 3.1
root2
└── X
`))
	})

	It("renders single root", func() {
		text := Render(rootnode2, DefaultVisitor, ts)
		Expect(text).To(Equal(`root2
└── X
`))
	})

	It("dereferences node values", func() {
		text := Render(&rootnode2, DefaultVisitor, ts)
		Expect(text).To(Equal(`root2
└── X
`))
		text = Render([]*Node{&rootnode2}, DefaultVisitor, ts)
		Expect(text).To(Equal(`root2
└── X
`))
	})

	It("renders roots maps", func() {
		text := Render(rootmap, sortingVisitor, ts)
		Expect(text).To(Equal(`alpharot
   • a
   • z
root1
│  • bar
│  • foo
├── 1
├── 2
│   ├── 2.1
│   │      • whoooosh
│   └── 2.2
└── 3
    └── 3.1
`))
	})

	It("renders map plainly", func() {
		text := RenderPlain(rootmap2)
		Expect(text).To(Equal(`root
|  * pr
` + "`" + `- 1
      * p1
      * p2
`))
	})

	It("renders maps", func() {
		text := Render(rootmap2, DefaultVisitor, ts)
		Expect(text).To(Equal(`root
│  • pr
└── 1
       • p1
       • p2
`))
	})

	It("renders fancy", func() {
		text := RenderFancy(rootmap2)
		Expect(strings.HasPrefix(text, "root\n")).To(BeTrue())
	})

	It("panics when rendering an unsupported roots type", func() {
		Expect(func() { Render(42, DefaultVisitor, ts) }).To(Panic())
		Expect(func() { Render([]int{42}, DefaultVisitor, ts) }).To(Panic())
	})

	It("panics when rendering incorrect node", func() {
		Expect(func() {
			// nolint structcheck
			type badNode struct {
				foo bool
			}
			Render(badNode{}, DefaultVisitor, ts)
		}).To(Panic())
		Expect(func() {
			type badNode struct {
				Foo bool `asciitree:"foo"`
			}
			Render(badNode{}, DefaultVisitor, ts)
		}).To(Panic())
	})

})
