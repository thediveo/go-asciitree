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
	"fmt"
	"iter"
	"reflect"
	"strings"
)

// renderSubtree returns an iterator that produces lines from recursively
// rendering the subtree starting at the passed tree node.
//
// The passed (tree) node can be any value, as long as the passed visitor is
// able to correctly determine the value's label as well as optional properties
// and children nodes.
//
// The styler parameter controls the output rendering, so a user-controllable
// style can be used while traversing the subtree.
func renderSubtree(node any, visitor Visitor, styler *TreeStyler) (lines iter.Seq[string]) {
	label, props, children := visitor.Get(node)
	return func(yield func(string) bool) {
		// produce the label of the passed node.
		if !yield(styler.renderNodeLabel(label)) {
			return
		}
		// next, produce the properties of this node.
		renderProp := styler.renderPropertyChildrenFollowing
		if len(children) == 0 {
			renderProp = styler.renderPropertyNoChildrenFollowing
		}
		for _, prop := range props {
			if !yield(renderProp(styler.renderProperty(prop))) {
				return
			}
		}
		// finally,f or each child subtree of the current tree node we first
		// render these subtrees and then indent the resulting text lines as
		// needed ... because we have to differentiate between intermediate
		// child nodes and the final child nodes in each subtree due to
		// different styling.
		last := len(children) - 1
		for idx := range len(children) {
			lines := renderSubtree(children[idx], visitor, styler)
			style := styler.renderBranchedNode
			styleButFirst := styler.indentLine
			if idx == last {
				style = styler.renderLastNode
				styleButFirst = styler.indentLineLastNode
			}
			for line := range lines {
				if !yield(style(line)) {
					return
				}
				style = styleButFirst
			}
		}
	}
}

// Render a tree (or a multi-root “tree” ... is that a forrest?) into a
// multi-line text string, using the supplied visitor and tree styler.
//
// The roots can be specified as a slice of structs, or also as a single struct.
// In every case, the passed root(s), as well as their subtree nodes need to
// have two struct fields exported and tagged as `asciitree:"label"` and
// `asciitree:"children"` respectively.
//
// For the visitor, you might want to simply use the DefaultVisitor that handles
// annotated structs and maps with well-known keys.
//
// As a styler, simply use DefaultTreeStyler, or the slightly more fancyful
// NewTreeStyler(LineStyle).
func Render(roots any, visitor Visitor, styler *TreeStyler) string {
	switch rv := reflect.Indirect(reflect.ValueOf(roots)); rv.Kind() {
	case reflect.Slice:
		// For a slice we need to iterate over all elements, passing the interface
		// of each element to the subtree renderer in turn. Please note that we
		// put the root element(s) first through the visitor just in case it wants
		// to sort nodes including root nodes.
		roots := visitor.Roots(roots)
		var result strings.Builder
		for idx := range len(roots) {
			for line := range renderSubtree(roots[idx], visitor, styler) {
				result.WriteString(line)
				result.WriteRune('\n')
			}
		}
		return result.String()
	case reflect.Struct:
		// A single root can be represented via a single struct for convenience,
		// so simply pass the struct value's interface to the subtree renderer,
		// and we're done.
		var result strings.Builder
		for line := range renderSubtree(roots, visitor, styler) {
			result.WriteString(line)
			result.WriteRune('\n')
		}
		return result.String()
	case reflect.Map:
		// A map with a "roots" key.
		maproots := rv.MapIndex(reflect.ValueOf("roots"))
		if maproots.Kind() == reflect.Invalid {
			var result strings.Builder
			for line := range renderSubtree(roots, visitor, styler) {
				result.WriteString(line)
				result.WriteRune('\n')
			}
			return result.String()
		}
		return Render(maproots.Interface(), visitor, styler)
	default:
		panic(fmt.Sprintf("unsupported roots value type: expected slice, map, or struct; got %T", roots))
	}
}

// RenderPlain renders a tree or multi-root tree into a multi-line text string
// using the default tree styling and (user data) visitor.
//
// The roots can be specified as a slice of structs, or also as a single
// struct. In every case, the passed root(s), as well as their subtree nodes
// need to have two struct fields exported and tagged as `asciitree:"label"`
// and `asciitree:"children"` respectively.
func RenderPlain(roots any) string {
	return Render(roots, DefaultVisitor, DefaultTreeStyler)
}

// RenderFancy works like RenderPlain, rendering a tree or multi-root tree
// into a multi-line text string, but it uses Unicode box characters to render
// the branch lines.
func RenderFancy(roots any) string {
	return Render(roots, DefaultVisitor, LineTreeStyler)
}
