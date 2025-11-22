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
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("visiting", func() {

	When("converting reflect.Value to any slices", func() {

		It("rejects non-slice reflect.Values", func() {
			Expect(anySlice(reflect.ValueOf(nil))).To(BeNil())
			Expect(anySlice(reflect.ValueOf(42))).To(BeNil())
			var fourtytwo = 42
			Expect(anySlice(reflect.ValueOf(&fourtytwo))).To(BeNil())
		})

		It("accepts a slice value", func() {
			Expect(anySlice(reflect.ValueOf([]string{"foo"}))).To(
				HaveExactElements("foo"))
		})

		It("accepts an interface value containing a slice", func() {
			foosl := []string{"foo"}
			ifaceV := reflect.New(reflect.TypeOf((*any)(nil)).Elem()).Elem()
			ifaceV.Set(reflect.ValueOf(foosl))
			Expect(anySlice(ifaceV)).To(
				HaveExactElements("foo"))
		})

	})

	When("working with node information", func() {

		Context("tagged structs", func() {

			It("gets nothing from nothing in Get() and Label()", func() {
				var testStruct struct{}

				label, properties, children := DefaultVisitor.Get(&testStruct)
				Expect(label).To(BeEmpty())
				Expect(properties).To(BeEmpty())
				Expect(children).To(BeEmpty())

				Expect(DefaultVisitor.Label(&testStruct)).To(BeEmpty())
			})

			It("picks up the correct fields in Get() and Label() ", func() {
				type T struct {
					Label    string   `asciitree:"label"`
					Props    []string `asciitree:"properties"`
					Children []T      `asciitree:"children"`
				}

				testTree := T{
					Label: "root",
					Props: []string{"fooprop", "barprop"},
					Children: []T{
						{Label: "child-B"},
						{Label: "child-A"},
					},
				}

				label, properties, children := NewMapStructVisitor(true, true).Get(testTree)
				Expect(label).To(Equal("root"))
				Expect(properties).To(HaveExactElements("barprop", "fooprop"))
				Expect(children).To(HaveExactElements(
					HaveField("Label", "child-A"),
					HaveField("Label", "child-B")))

				Expect(DefaultVisitor.Label(testTree)).To(Equal("root"))
			})

		})

		Context("maps with well-known keys", func() {

			It("gets nothing from nothing in Get() and Label()", func() {
				var testTree map[string]string

				label, properties, children := DefaultVisitor.Get(testTree)
				Expect(label).To(BeEmpty())
				Expect(properties).To(BeEmpty())
				Expect(children).To(BeEmpty())

				Expect(DefaultVisitor.Label(testTree)).To(BeEmpty())
			})

			It("picks up the correct values in Get() and Label()", func() {
				type M map[string]any
				testTree := M{
					"label":      "root",
					"properties": []string{"fooprop", "barprop"},
					"children": []M{
						{"label": "child-B"},
						{"label": "child-A"},
					},
				}

				label, properties, children := (&MapStructVisitor{
					SortNodes:      true,
					SortProperties: true,
				}).Get(testTree)
				Expect(label).To(Equal("root"))
				Expect(properties).To(HaveExactElements("barprop", "fooprop"))
				Expect(children).To(HaveExactElements(
					HaveKeyWithValue("label", "child-A"),
					HaveKeyWithValue("label", "child-B")))

				Expect(DefaultVisitor.Label(testTree)).To(Equal("root"))
			})
		})

	})

	When("handling structs", func() {

		type T struct {
			bar        int    // nolint
			MyLabel    string `asciitree:"label"`
			MyChildren []T    `asciitree:"children"`
		}

		testTree := T{
			MyLabel: "root",
			MyChildren: []T{
				{MyLabel: "2"},
				{MyLabel: "1"},
			},
		}

		It("uses fields with well-known tags in Get", func() {
			label, _, children := DefaultVisitor.Get(&testTree)
			Expect(label).To(Equal(testTree.MyLabel))
			Expect(children).To(HaveLen(2))
		})

		It("uses a field with a well-known tag in Label", func() {
			Expect(DefaultVisitor.Label(&testTree)).To(Equal(testTree.MyLabel))
		})

	})

	When("handling maps", func() {

		type M map[string]any

		var testMap = M{
			"label":      "root",
			"properties": []string{"someprop: value"},
			"children": []M{
				{"label": "child 3"},
				{"label": "child 2", "properties": []string{"foo", "bar"}},
				{"label": "child 1"},
			},
		}

		It("uses well-known keys in Get", func() {
			label, props, children := DefaultVisitor.Get(testMap)
			Expect(label).To(Equal(testMap["label"]))
			Expect(props).To(HaveExactElements("someprop: value"))
			Expect(children).To(HaveLen(3))
		})

		It("uses a well-known key in Label", func() {
			Expect(DefaultVisitor.Label(testMap)).To(Equal(testMap["label"]))
		})

		It("retrieves 'children' nodes and sorts them by their 'label's", func() {
			_, _, children := DefaultVisitor.Get(testMap)
			sortedChildren := DefaultVisitor.sortedNodes(children)
			Expect(sortedChildren).ToNot(Equal(children))
			Expect(sortedChildren).To(HaveLen(3))
			label := DefaultVisitor.Label(sortedChildren[0])
			Expect(label).To(Equal("child 1"))
		})

	})

	It("panics when presented with neither struct nor map", func() {
		Expect(func() { _, _, _ = DefaultVisitor.nodeDetails(42) }).To(
			PanicWith(MatchRegexp(`unsupported asciitree node.*type int`)))
		Expect(func() { _ = DefaultVisitor.nodeLabel(42) }).To(
			PanicWith(MatchRegexp(`unsupported asciitree node.*type int`)))
	})

	When("working with roots", func() {

		It("unwraps a slice of roots", func() {
			type T struct {
				Label string `asciitree:"label"`
			}

			trees := []T{
				{Label: "ruth"},
				{Label: "root"},
			}

			Expect(DefaultVisitor.Roots(trees)).To(
				HaveExactElements(
					HaveField("Label", "ruth"),
					HaveField("Label", "root")))
			Expect((&MapStructVisitor{SortNodes: true}).Roots(trees)).To(
				HaveExactElements(
					HaveField("Label", "root"),
					HaveField("Label", "ruth")))
		})

		Context("handling structs", func() {

			It("handles a root-less struct as a single root", func() {
				type T struct {
					Label string `asciitree:"label"`
				}

				tree := T{
					Label: "ruth",
				}

				Expect(DefaultVisitor.Roots(tree)).To(
					HaveExactElements(HaveField("Label", "ruth")))
			})

			It("picks up the roots field", func() {
				type T struct {
					Label string
				}
				type R struct {
					Roots []T `asciitree:"roots"`
				}

				roots := R{
					Roots: []T{
						{Label: "ruth"},
						{Label: "root"},
					},
				}
				Expect(DefaultVisitor.Roots(roots)).To(
					HaveExactElements(
						HaveField("Label", "ruth"),
						HaveField("Label", "root")))
			})

		})

		Context("handling maps", func() {

			It("handles roots-less map as itself", func() {
				type T map[string]any

				tree := T{
					"label": "foo",
				}

				Expect(DefaultVisitor.Roots(tree)).To(HaveExactElements(
					Equal(tree)))
			})

			It("handles a map with a roots slice value", func() {
				type T map[string]any

				roots := T{
					"roots": []T{
						{"label": "ruth"},
						{"label": "root"},
					},
				}

				Expect(DefaultVisitor.Roots(roots)).To(HaveExactElements(
					HaveKeyWithValue("label", "ruth"),
					HaveKeyWithValue("label", "root")))
			})

			It("handles a map with a roots slice value", func() {
				type T map[string]any

				roots := T{
					"roots": T{"label": "ruth"},
				}

				Expect(DefaultVisitor.Roots(roots)).To(HaveExactElements(
					HaveKeyWithValue("label", "ruth")))
			})

		})

		It("panics when presented with anything else than slice, struct, map", func() {
			Expect(func() { _ = DefaultVisitor.Roots(42) }).To(
				PanicWith(MatchRegexp(`but got int`)))
		})

	})

	/*
		// nolint structcheck
		type S struct {
			bar int
			Foo string `asciitree:"label"`
			Baz []S    `asciitree:"children"`
		}

		tree := S{
			Foo: "root",
			Baz: []S{
				{Foo: "2"},
				{Foo: "1"},
			},
		}

		// nolint structcheck
		type RS struct {
			fake  bool `asciitree:"label"`
			Rootz []S  `asciitree:"roots"`
		}

		roots := RS{
			Rootz: []S{tree, tree.Baz[0]},
		}
		_ = roots

		type M map[string]any

		var mapp = M{
			"label":      "root",
			"properties": []string{"someprop: value"},
			"children": []M{
				{"label": "child 3"},
				{"label": "child 2", "properties": []string{"foo", "bar"}},
				{"label": "child 1"},
			},
		}

		})
	*/

	/*
		Describe("roots", func() {

			d := &MapStructVisitor{}

			It("panics on unsupported roots", func() {
				Expect(func() { d.Roots(reflect.ValueOf(42)) }).To(Panic())
			})

			Describe("maps of roots", func() {
				It("traverses pointer to map root", func() {
					r := d.Roots(reflect.ValueOf(&mapp))
					Expect(len(r)).To(Equal(1))
					Expect(r[0].Interface().(M)).To(Equal(mapp))
				})

				It("traverses slice of map roots", func() {
					r := d.Roots(reflect.ValueOf([]M{mapp}))
					Expect(len(r)).To(Equal(1))
					Expect(r[0].Interface().(M)).To(Equal(mapp))
				})

				It("traverses map root", func() {
					r := d.Roots(reflect.ValueOf(mapp))
					Expect(len(r)).To(Equal(1))
					Expect(r[0].Interface().(M)).To(Equal(mapp))
				})

				It("traverses map roots scalar", func() {
					r := d.Roots(reflect.ValueOf(M{"roots": mapp}))
					Expect(len(r)).To(Equal(1))
					Expect(r[0].Interface().(M)).To(Equal(mapp))
				})

				It("traverses map roots slice", func() {
					// well, so much for the easy golang expressioness ... it's so
					// much better than C's brackets-and-stars wilderness! So MUCH!!!
					r := d.Roots(reflect.ValueOf(M{"roots": []M{mapp, mapp["children"].([]M)[0]}}))
					Expect(len(r)).To(Equal(2))
					Expect(r[0].Interface().(M)).To(Equal(mapp))
					Expect(r[1].Interface().(M)).To(Equal(mapp["children"].([]M)[0]))
				})
			})

			Describe("structs of roots", func() {
				It("panics on unknown root asciitree tags", func() {
					// nolint structcheck
					type RS struct {
						fake bool `asciitree:"foobar"`
					}
					Expect(func() { d.Roots(reflect.ValueOf(RS{})) }).To(Panic())
				})

				It("traverses struct root", func() {
					r := d.Roots(reflect.ValueOf(roots))
					Expect(len(r)).To(Equal(2))
					Expect(r[0].Interface().(S)).To(Equal(roots.Rootz[0]))
					Expect(r[1].Interface().(S)).To(Equal(roots.Rootz[1]))
				})

				It("traverses single root struct", func() {
					r := d.Roots(reflect.ValueOf(tree))
					Expect(len(r)).To(Equal(1))
				})
			})

		})
	*/

})
