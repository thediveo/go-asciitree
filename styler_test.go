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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Styler", func() {

	Describe("helpers", func() {
		Context("repeat()", func() {
			It("doesn't panic when repeating a string less than zero times", func() {
				Expect(func() { repeat("foo", -1) }).ToNot(Panic()) // nolint staticcheck
			})

			It("returns an empty string when not repeating", func() {
				Expect(repeat("foo", -1)).To(Equal(""))
			})
		})
	})

	Describe("rendering", func() {
		It("gives correct text output", func() {
			s := NewTreeStyler(ASCIIStyle)
			s.ChildIndent = 4
			s.PropIndent = 4
			Expect(s.renderBranchedNode("foo")).To(Equal("+-- foo"))
			Expect(s.renderLastNode("foo")).To(Equal("`-- foo"))
			Expect(s.indentLine(s.renderBranchedNode("foo"))).To(Equal("|   +-- foo"))
			Expect(s.renderPropertyNoChildrenFollowing("proo")).To(Equal("    * proo"))
			Expect(s.renderPropertyNoChildrenFollowing("proo")).To(Equal("    * proo"))
		})
	})

})
