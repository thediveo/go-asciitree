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
	"reflect"
	"strings"
)

// Recursively renders a subtree starting at a specific tree node, returning
// an array of the resulting text lines. The tree node can be any struct, as
// long as it has the required "asciitree" tags on the relevant exported(!)
// fields.
//
// The styler parameter controls the output rendering, so a user-controllable
// style can be used while traversing the subtree.
//
func renderSubtree(node interface{}, visitor Visitor, styler *TreeStyler) (lines []string) {
	label, props, children := visitor.Get(reflect.ValueOf(node))

	lines = append(lines, styler.renderNodeLabel(label))

	// Render properties of this node.
	renderProp := styler.renderPropertyChildrenFollowing
	if children.Len() == 0 {
		renderProp = styler.renderPropertyNoChildrenFollowing
	}
	for _, prop := range props {
		line := renderProp(styler.renderProperty(prop))
		lines = append(lines, line)
	}

	// For each child subtree of the current tree node we first
	// render these subtrees and then indent the resulting text
	// lines as needed ... because we have to differentiate between
	// intermediate child nodes and the final child nodes in each
	// subtree due to different styling.
	last := children.Len() - 1
	for idx := 0; idx <= last; idx++ {
		childLines := renderSubtree(
			reflect.Indirect(children.Index(idx)).Interface(), visitor, styler)
		if idx != last {
			lines = append(lines, styler.renderBranchedNode(childLines[0]))
			for _, childLine := range childLines[1:] {
				lines = append(lines, styler.indentLine(childLine))
			}
		} else {
			lines = append(lines, styler.renderLastNode(childLines[0]))
			for _, childLine := range childLines[1:] {
				lines = append(lines, styler.indentLineLastNode(childLine))
			}
		}
	}
	return
}

// Render renders a tree (or multi-root tree) into a multi-line text string,
// using the supplied visitor and tree styler.
//
// The roots can be specified as a slice of structs, or also as a single
// struct. In every case, the passed root(s), as well as their subtree nodes
// need to have two struct fields exported and tagged as `asciitree:"label"`
// and `asciitree:"children"` respectively.
//
// For the visitor, you might want to simply use the DefaultVisitor that
// handles annotated structs and maps with well-known keys.
//
// As a styler, simply use DefaultTreeStyler, or the slightly more fancyful
// NewTreeStyler(LineStyle).
func Render(roots interface{}, visitor Visitor, styler *TreeStyler) string {
	rv := reflect.Indirect(reflect.ValueOf(roots))
	switch rv.Kind() {
	// For a slice we need to iterate over all elements, passing the interface
	// of each element to the subtree renderer in turn. Please note that we
	// put the root element(s) first through the visitor just in case it wants
	// to sort nodes including root nodes.
	case reflect.Slice:
		roots := visitor.Roots(rv)
		result := ""
		count := len(roots)
		for idx := 0; idx < count; idx++ {
			lines := renderSubtree(roots[idx].Interface(), visitor, styler)
			result = result + strings.Join(lines, "\n") + "\n"
		}
		return result
	// A single root can be represented via a single struct for convenience,
	// so simply pass the struct value's interface to the subtree renderer,
	// and we're done.
	case reflect.Struct:
		lines := renderSubtree(rv.Interface(), visitor, styler)
		return strings.Join(lines, "\n") + "\n"
	//
	case reflect.Map:
		maproots := rv.MapIndex(reflect.ValueOf("roots"))
		if maproots.Kind() == reflect.Invalid {
			lines := renderSubtree(roots, visitor, styler)
			return strings.Join(lines, "\n") + "\n"
		}
		return Render(maproots.Interface(), visitor, styler)
	default:
		panic(fmt.Sprintf("unsupported roots type %q", rv.Kind()))
	}
}

// RenderPlain renders a tree or multi-root tree into a multi-line text string
// using the default tree styling and (user data) visitor.
//
// The roots can be specified as a slice of structs, or also as a single
// struct. In every case, the passed root(s), as well as their subtree nodes
// need to have two struct fields exported and tagged as `asciitree:"label"`
// and `asciitree:"children"` respectively.
func RenderPlain(roots interface{}) string {
	return Render(roots, DefaultVisitor, DefaultTreeStyler)
}

// RenderFancy works like RenderPlain, rendering a tree or multi-root tree
// into a multi-line text string, but it uses Unicode box characters to render
// the branch lines.
func RenderFancy(roots interface{}) string {
	return Render(roots, DefaultVisitor, LineTreeStyler)
}
