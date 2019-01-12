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
)

// TreeStyle defines the ASCII art elements required for "painting" beautiful
// ASCII trees. In this case, we also subsume Unicode under ASCII for a suitable
// definition of "ASCII" and "Unicode".
type TreeStyle struct {
	Fork     string // depicts forking off a node, such as "├".
	Nodeconn string // depicts the branch going to a node "leaf", such as "─".
	Nofork   string // depicts a continuing vertical main branch, such as "│".
	Lastnode string // depicts a vertical main branch ending in a node, such as "└".
	Property string // depicts a property, such as "•"
}

// ASCIIStyle styles rendered trees using only "pure" ASCII characters,
// without any help of special line or box characters.
var ASCIIStyle = TreeStyle{
	Fork:     "+",
	Nodeconn: "-",
	Nofork:   "|",
	Lastnode: "`",
	Property: "*",
}

// LineStyle styles ASCII trees using Unicode line characters.
var LineStyle = TreeStyle{
	Fork:     "├", // Don't print this on an FX-80/100 ;)
	Nodeconn: "─",
	Nofork:   "│",
	Lastnode: "└",
	Property: "•",
}

// TreeStyler describes the tree branch and node properties indentations, as
// well as the style of "line art" to use when rendering ASCII trees.
type TreeStyler struct {
	Style       TreeStyle // The specific TreeStyle to use, such as ASCIIStyle, or LineStyle.
	ChildIndent int       // The indentation of child nodes.
	PropIndent  int       // The indentation of properties w.r.t. their node
}

// DefaultTreeStyler offers a pure ASCII tree styler, using only "safe"
// ASCII characters, but no Unicode characters. Ideal for the lovers of
// unwatered ASCII art.
var DefaultTreeStyler = NewTreeStyler(ASCIIStyle)

// LineTreeStyler uses Unicode box characters to draw slightly fancy tree
// branches.
var LineTreeStyler = NewTreeStyler(LineStyle)

// NewTreeStyler returns a Style object suitable for use with rendering trees.
func NewTreeStyler(style TreeStyle) *TreeStyler {
	s := new(TreeStyler)
	s.Style = style
	s.ChildIndent = 3
	s.PropIndent = 3
	return s
}

// Defaults to no adornments to node labels
func (s *TreeStyler) renderNodeLabel(label string) string {
	return label
}

func (s *TreeStyler) renderBranchedNode(label string) string {
	return s.Style.Fork +
		repeat(s.Style.Nodeconn, s.ChildIndent-2) +
		" " +
		label
}

func (s *TreeStyler) renderLastNode(label string) string {
	return s.Style.Lastnode +
		repeat(s.Style.Nodeconn, s.ChildIndent-2) +
		" " +
		label
}

func (s *TreeStyler) indentLine(line string) string {
	return s.Style.Nofork +
		repeat(" ", s.ChildIndent-2) +
		" " +
		line
}

func (s *TreeStyler) indentLineLastNode(line string) string {
	return repeat(" ", s.ChildIndent) + line
}

func (s *TreeStyler) renderProperty(prop string) string {
	return prop
}

func (s *TreeStyler) renderPropertyNoChildrenFollowing(prop string) string {
	return repeat(" ", s.PropIndent) +
		s.Style.Property +
		" " +
		prop
}

func (s *TreeStyler) renderPropertyChildrenFollowing(prop string) string {
	return s.Style.Nofork +
		repeat(" ", s.PropIndent-1) +
		s.Style.Property +
		" " +
		prop
}

// Like strings.Repeat, but without its panic if count is less than
// zero.
func repeat(s string, count int) string {
	if count > 0 {
		return strings.Repeat(s, count)
	}
	return ""
}
