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
	"slices"
	"sort"
	"strings"
)

// Visitor looks into user-defined data structures, trying to locate the
// anointed, erm, annotated (“tagged”) user-data entries which represent node
// labels, properties, and children. The Visitor interface is used by tree
// renderers for visiting nodes in order to retrieve their tree-relevant
// information while traversing trees.
//
// Visitors thus avoid having to implement an asciitree-specific interface on
// user data-structures, that is, they avoid adaptors. In many use cases, a
// MapStructVisitor will suffice as it implements the visitor pattern on
// tagged structs and maps with well-known keys.
//
// Please note that Visitors do not traverse; that is the job of the asciitree
// Render...() functions.
type Visitor interface {
	Roots(roots any) (children []any)
	Label(node any) (label string)
	Get(node any) (label string, properties []string, children []any)
}

// DefaultVisitor provides visitor capable of traversing (annotated) maps and
// structs.
var DefaultVisitor = &MapStructVisitor{}

// MapStructVisitor visits tagged ("annotated") user-defined structs as well
// as maps (the latter using well-known keys) and retrieves their
// tree-relevant data. For convenience, it also handles slices and pointers to
// structs and maps.
type MapStructVisitor struct {
	Visitor
	SortNodes      bool
	SortProperties bool
}

var _ Visitor = (*MapStructVisitor)(nil)

// NewMapStructVisitor creates a visitor that optionally sorts nodes and their
// properties.
func NewMapStructVisitor(sortNodes bool, sortProperties bool) *MapStructVisitor {
	return &MapStructVisitor{SortNodes: sortNodes, SortProperties: sortProperties}
}

// Roots returns the list of root nodes, while handling different types of
// Roots data types; for instance, struct, []struct, map, and []map, as well
// as pointers.
func (v *MapStructVisitor) Roots(roots any) []any {
	switch rv := reflect.Indirect(reflect.ValueOf(roots)); rv.Kind() {
	case reflect.Slice:
		// For a slice we need to iterate over all elements, so we return all
		// elements as the list of "children". If this visitor is configured to
		// sort by label, then we also need to sort the roots.
		roots := anySlice(rv)
		if !v.SortNodes {
			return roots
		}
		return v.sortedNodes(roots)
	case reflect.Struct:
		// A single root can be represented via a single struct for convenience,
		// so simply return a list of "children" consisting only of this single
		// struct itself.
		si := structFieldInfo(rv)
		if si.RootsPath == nil {
			return []any{roots}
		}
		return v.Roots(rv.FieldByIndex(si.RootsPath).Interface())
	case reflect.Map:
		// Finally, roots can also be stored in a map using a well-known key
		// named "roots". If that key is present, then it must be a list of
		// children, otherwise return a list of children consisting only if this
		// map itself because it's already a child.
		maproots := rv.MapIndex(reflect.ValueOf("roots"))
		switch maproots.Kind() {
		case reflect.Invalid:
			// Nope, no such "roots" key, so the root given is the only one root
			// node itself, so return it as a list of exactly one root node.
			return []any{roots}
		default:
			// The roots element may be represented by map or struct, or it might
			// be a slice of them; especially due to the latter case we have to
			// create a Value slice with the individual elements from the "roots"
			// property. Unfortunately, we cannot simply return the slice Value
			// itself, but instead need to create a new slice of Values
			// referencing the elements of the original slice.
			elem := maproots.Elem()
			if elem.Kind() == reflect.Slice {
				return v.Roots(elem.Interface())
			}
			return []any{reflect.Indirect(maproots).Interface()}
		}
	default:
		panic(fmt.Sprintf("expecting roots to be a slice, struct, or map, but got %T", roots))
	}
}

// Label returns the label for a tree node.
func (v *MapStructVisitor) Label(node any) (label string) {
	return v.nodeLabel(node)
}

// Get returns the label, properties, and children of a tree node, hiding
// pesty details about how to fetch them from tagged structs or maps with
// well-known fields.
func (v *MapStructVisitor) Get(node any) (label string, properties []string, children []any) {
	label, properties, children = v.nodeDetails(node)
	if v.SortProperties {
		properties = slices.Clone(properties)
		sort.Strings(properties)
	}
	return label, properties, children
}

func (v *MapStructVisitor) nodeLabel(node any) string {
	switch node := reflect.Indirect(reflect.ValueOf(node)); node.Kind() {
	case reflect.Struct:
		si := structFieldInfo(node)
		if si.LabelPath == nil {
			return ""
		}
		return node.FieldByIndex(si.LabelPath).String()
	case reflect.Map:
		labelV := node.MapIndex(reflect.ValueOf("label"))
		if labelV.Kind() == reflect.Interface {
			labelV = labelV.Elem()
		}
		if labelV.Kind() != reflect.String {
			return ""
		}
		return labelV.Interface().(string)
	default:
		panic(fmt.Sprintf("unsupported asciitree node or root type %T", node.Interface()))
	}
}

// Internal helper to either only retrieve the label for a node, or the label,
// properties, and children. Please note that we don't sort here; this is
// really only the helper for retrieving.
func (v *MapStructVisitor) nodeDetails(node any) (label string, properties []string, children []any) {
	switch node := reflect.Indirect(reflect.ValueOf(node)); node.Kind() {
	case reflect.Struct:
		// Grab the values for a node label, its properties, and its children,
		// if there are fields known to have them – based on their field
		// tags.
		si := structFieldInfo(node)
		if si.LabelPath != nil {
			label = node.FieldByIndex(si.LabelPath).String()
		}
		if si.PropertiesPath != nil {
			properties = node.FieldByIndex(si.PropertiesPath).Interface().([]string)
		}
		if si.ChildrenPath == nil {
			return
		}
		children = anySlice(node.FieldByIndex(si.ChildrenPath))
		if !v.SortNodes {
			return
		}
		children = v.sortedNodes(children)
		return
	case reflect.Map:
		// Gets the (well-known) key-values for label, properties, and children in
		// a map. Again, all these keys-values are optional and will default to
		// zero if missing.
		if lbl := node.MapIndex(reflect.ValueOf("label")); lbl.Kind() != reflect.Invalid {
			label, _ = lbl.Interface().(string)
		}
		if pps := node.MapIndex(reflect.ValueOf("properties")); pps.Kind() != reflect.Invalid {
			properties, _ = pps.Elem().Interface().([]string)
		}
		if chs := node.MapIndex(reflect.ValueOf("children")); chs.Kind() != reflect.Invalid {
			children = anySlice(chs)
			if v.SortNodes {
				children = v.sortedNodes(children)
			}
		}
		return
	default:
		panic(fmt.Sprintf("unsupported asciitree node or root type %T", node.Interface()))
	}
}

// sortedNodes returns a new slice of sorted nodes from the passed slice of
// nodes, sorted by lexicographically by their labels.
func (v *MapStructVisitor) sortedNodes(nodes []any) []any {
	type labelledNode struct {
		Label string
		Node  any
	}
	l := len(nodes)
	labelledNodes := make([]labelledNode, l)
	for idx := range nodes {
		labelledNodes[idx] = labelledNode{Label: v.Label(nodes[idx]), Node: nodes[idx]}
	}
	slices.SortStableFunc(labelledNodes, func(a, b labelledNode) int {
		return strings.Compare(a.Label, b.Label)
	})
	sortednodes := make([]any, l)
	for idx := range l {
		sortednodes[idx] = labelledNodes[idx].Node
	}
	return sortednodes
}

// anySlice returns an []any value whose elements are the slice elements
// contained in the passed reflect.Value (unpacking an interface value where
// necessary), or nil if the passed reflect.Value is not a slice.
func anySlice(v reflect.Value) []any {
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice {
		return nil
	}
	l := v.Len()
	anyslice := make([]any, l)
	for idx := range l {
		anyslice[idx] = v.Index(idx).Interface()
	}
	return anyslice
}
