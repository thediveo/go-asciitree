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

var _ = Describe("field index cache", func() {

	It("validates asciitree tags", func() {
		vals, ok := asciitreeTagValues("")
		Expect(ok).To(BeTrue())
		Expect(vals).To(BeEmpty())

		vals, ok = asciitreeTagValues(",   ,, ")
		Expect(ok).To(BeTrue())
		Expect(vals).To(BeEmpty())

		vals, ok = asciitreeTagValues(" label")
		Expect(ok).To(BeTrue())
		Expect(vals).To(Equal([]string{"label"}))

		vals, ok = asciitreeTagValues("roots,label,properties,children")
		Expect(ok).To(BeTrue())
		Expect(vals).To(Equal([]string{"roots", "label", "properties", "children"}))

		vals, ok = asciitreeTagValues("label,foo,bar")
		Expect(ok).To(BeFalse())
		Expect(vals).To(Equal([]string{"foo", "bar"}))
	})

	It("finds struct field indices (even embedded anonymous ones)", func() {
		// nolint structcheck
		type X struct {
			bar int
		}

		si := structInfo(reflect.ValueOf(42))
		Expect(si).To(BeNil())

		fieldCache = make(map[reflect.Type]*fieldCacheItem)
		si = structInfo(reflect.ValueOf(X{}))
		Expect(len(fieldCache)).To(Equal(1))
		Expect(*si).To(Equal(fieldCacheItem{-1, -1, -1, -1}))

		type S struct {
			X
			Foo  string   `asciitree:"label"`
			Poop []string `asciitree:"properties"`
			Baz  []S      `asciitree:"children"`
			Rotz []S      `asciitree:"roots"`
		}
		si = structInfo(reflect.ValueOf(S{}))
		Expect(*si).To(Equal(fieldCacheItem{1, 2, 3, 4}))

		type SKO struct {
			S
			Rootz []S `asciitree:"roots"`
		}
		Expect(func() { structInfo(reflect.ValueOf(SKO{})) }).To(Panic())

		type SOK struct {
			S     S
			Rootz []S `asciitree:"roots"`
		}
		Expect(func() { si = structInfo(reflect.ValueOf(SOK{})) }).ToNot(Panic())
		Expect(*si).To(Equal(fieldCacheItem{-1, -1, -1, 1}))

		type SKO2 struct {
			S
			Ohno bool `asciitree:"label"`
		}
		Expect(func() { structInfo(reflect.ValueOf(SKO2{})) }).To(Panic())

		type SKO3 struct {
			S
			Ohno bool `asciitree:"properties"`
		}
		Expect(func() { structInfo(reflect.ValueOf(SKO3{})) }).To(Panic())

		type SKO4 struct {
			S
			Ohno bool `asciitree:"children"`
		}
		Expect(func() { structInfo(reflect.ValueOf(SKO4{})) }).To(Panic())

		// nolint structcheck
		type TL struct {
			laberl string `asciitree:"label"`
		}
		type SKO666 struct {
			S
			TL
		}
		Expect(func() { structInfo(reflect.ValueOf(SKO666{})) }).To(Panic())

		// nolint structcheck
		type TL2 struct {
			prups string `asciitree:"properties"`
		}
		type SKO667 struct {
			S
			TL2
		}
		Expect(func() { structInfo(reflect.ValueOf(SKO667{})) }).To(Panic())

		// nolint structcheck
		type TL3 struct {
			kinners string `asciitree:"children"`
		}
		type SKO668 struct {
			S
			TL3
		}
		Expect(func() { structInfo(reflect.ValueOf(SKO668{})) }).To(Panic())

		// nolint structcheck
		type TL4 struct {
			absoluterrotz string `asciitree:"roots"`
		}
		type SKO669 struct {
			S
			TL4
		}
		Expect(func() { structInfo(reflect.ValueOf(SKO669{})) }).To(Panic())

	})

})
