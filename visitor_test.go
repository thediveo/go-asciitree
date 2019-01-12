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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("visiting", func() {

	type S struct {
		bar int
		Foo string `asciitree:"label"`
		Baz []S    `asciitree:"children"`
	}

	tree := S{
		Foo: "root",
		Baz: []S{
			S{Foo: "2"},
			S{Foo: "1"},
		},
	}

	type RS struct {
		fake  bool `asciitree:"label"`
		Rootz []S  `asciitree:"roots"`
	}

	roots := RS{
		Rootz: []S{tree, tree.Baz[0]},
	}

	type M map[string]interface{}

	var mapp M = M{
		"label":      "root",
		"properties": []string{"someprop: value"},
		"children": []M{
			M{"label": "child 3"},
			M{"label": "child 2", "properties": []string{"foo", "bar"}},
			M{"label": "child 1"},
		},
	}

	It("sorts labelled nodes", func() {
		_, _, children := DefaultVisitor.Get(reflect.ValueOf(mapp))
		label := DefaultVisitor.Label(children.Index(0))
		cs := sortNodes(DefaultVisitor, children, true)
		Expect(cs).ToNot(Equal(children))
		Expect(cs.Len()).To(Equal(3))
		label = DefaultVisitor.Label(cs.Index(0))
		Expect(label).To(Equal("child 1"))
	})

	Describe("visits", func() {

		It("struct", func() {
			label, _, children := DefaultVisitor.Get(reflect.ValueOf(&tree))
			Expect(label).To(Equal("root"))
			Expect(children.Len()).To(Equal(2))

			label = DefaultVisitor.Label(reflect.ValueOf(&tree))
			Expect(label).To(Equal("root"))
		})

		It("map", func() {
			label, props, children := DefaultVisitor.Get(reflect.ValueOf(mapp))
			Expect(label).To(Equal("root"))
			Expect(len(props)).To(Equal(1))
			Expect(children.Len()).To(Equal(3))

			label, props, _ = DefaultVisitor.Get(children.Index(1))
			Expect(label).To(Equal("child 2"))
			Expect(len(props)).To(Equal(2))
		})

	})

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

})
