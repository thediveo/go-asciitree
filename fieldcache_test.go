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
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("field index cache", func() {

	DescribeTable("checking asciitree tags",
		func(tag, value string, expected bool) {
			field := reflect.StructField{
				Tag: reflect.StructTag(tag),
			}
			Expect(hasAsciitreeTagValue(field, value)).To(Equal(expected))
		},
		Entry(nil, `foo:"bar"`, "", false),
		Entry(nil, `foo:"bar" asciitree:"label"`, "label", true),
		Entry(nil, `foo:"bar" asciitree:"foo"`, "label", false),
	)

	When("looking up types in the cache", func() {

		var cache *sync.Map

		BeforeEach(func() {
			cache = new(sync.Map)
		})

		It("returns nil for non-struct values", func() {
			Expect(structInfoCache(cache, reflect.ValueOf(42))).To(BeNil())
		})

		It("shrugs when there are no magic fields", func() {
			type T struct {
				foo int
			}
			si := structInfoCache(cache, reflect.ValueOf(T{foo: 42}))
			Expect(si).To(And(
				HaveField("LabelPath", BeNil()),
				HaveField("PropertiesPath", BeNil()),
				HaveField("ChildrenPath", BeNil()),
				HaveField("RootsPath", BeNil())))
			siAgain := structInfoCache(cache, reflect.ValueOf(T{}))
			Expect(siAgain).To(BeIdenticalTo(si))
		})

		It("finds magic fields", func() {
			type T struct {
				Foo string `asciitree:"label"`
			}
			type U struct {
				Bar int
				Baz []string `asciitree:"properties"`
				T
				Coolz []U `asciitree:"children"`
				Ruhtz []T `asciitree:"roots"`
			}
			Expect(structInfoCache(cache, reflect.ValueOf(U{}))).To(And(
				HaveField("LabelPath", HaveExactElements(2, 0)),
				HaveField("PropertiesPath", HaveExactElements(1)),
				HaveField("ChildrenPath", HaveExactElements(3)),
				HaveField("RootsPath", HaveExactElements(4))))
		})

	})

})
