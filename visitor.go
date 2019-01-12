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
	"sort"
)

// Visitor looks into user-defined data structures, trying to locate the
// anointed, erm, annotated ("tagged") user-data entries which represent node
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
	Roots(roots reflect.Value) (children []reflect.Value)
	Label(node reflect.Value) (label string)
	Get(node reflect.Value) (label string, properties []string, children reflect.Value)
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

// NewMapStructVisitor creates a visitor that optionally sorts nodes and their
// properties.
func NewMapStructVisitor(sortNodes bool, sortProperties bool) *MapStructVisitor {
	return &MapStructVisitor{SortNodes: sortNodes, SortProperties: sortProperties}
}

// Roots returns the list of root nodes, while handling different types of
// Roots data types; for instance, struct, []struct, map, and []map, as well
// as pointers.
func (v *MapStructVisitor) Roots(roots reflect.Value) (children []reflect.Value) {
	roots = reflect.Indirect(roots)
	switch roots.Kind() {
	// For a slice we need to iterate over all elements, so we return all
	// elements as the list of "children". If this visitor is configured to
	// sort by label, then we also need to sort the roots; we can do the
	// sorting in place, as we already had to create a shallow slice copy of
	// root/children anyway.
	case reflect.Slice:
		count := roots.Len()
		children = make([]reflect.Value, count)
		for idx := 0; idx < count; idx++ {
			children[idx] = roots.Index(idx)
		}
		if v.SortNodes {
			sortNodes(v, reflect.ValueOf(children), false)
		}
		return
	// A single root can be represented via a single struct for convenience,
	// so simply return a list of "children" consisting only of this single
	// struct itself.
	case reflect.Struct:
		fci := structInfo(roots)
		if fci.rootsIndex >= 0 {
			return v.Roots(roots.Field(fci.rootsIndex))
		}
		return []reflect.Value{roots}
	// Moreover, roots can also be stored in a map using a well-known key
	// named "roots". If that key is present, then it must be a list of
	// children, otherwise return a list of children consisting only if
	// this map itself because it's already a child.
	case reflect.Map:
		maproots := roots.MapIndex(reflect.ValueOf("roots"))
		switch maproots.Kind() {
		// Nope, no such "roots" key, so the root given is the only one root
		// node itself, so return it as a list of exactly one root node.
		case reflect.Invalid:
			return []reflect.Value{roots}
		// The roots element may be represented by map or struct, or it might
		// be a slice of them; especially due to the latter case we have to
		// create a Value slice with the individual elements from the "roots"
		// property. Unfortunately, we cannot simply return the slice Value
		// itself, but instead need to create a new slice of Values
		// referencing the elements of the original slice.
		default:
			elem := maproots.Elem()
			if elem.Kind() == reflect.Slice {
				return v.Roots(elem)
			}
			return []reflect.Value{reflect.Indirect(maproots)}
		}
	default:
		panic(fmt.Sprintf("unsupported roots type %q", roots.Kind()))
	}
}

// Label returns the label for a tree node.
func (v *MapStructVisitor) Label(node reflect.Value) (label string) {
	label, _, _ = getNodeData(v, node, true)
	return
}

// Get returns the label, properties, and children of a tree node, hiding
// pesty details about how to fetch them from tagged structs or maps with
// well-known fields.
func (v *MapStructVisitor) Get(node reflect.Value) (label string, properties []string, children reflect.Value) {
	label, properties, children = getNodeData(v, node, false)
	if v.SortProperties {
		sort.Strings(properties)
	}
	return
}

// Internal helper to either only retrieve the label for a node, or the label,
// properties, and children. Please note that we don't sort here; this is
// really only the helper for retrieving.
func getNodeData(v *MapStructVisitor, node reflect.Value, labelOnly bool) (label string, properties []string, children reflect.Value) {
	node = reflect.Indirect(node)
	switch node.Kind() {
	// Gets the fields for label, properties, and children in a struct; all
	// these fields are optional and will default to zeroed values if missing.
	case reflect.Struct:
		// Grab the values for a node label, its properties, and its children,
		// if there are fields known to have them -- based on their field
		// tags.
		fci := structInfo(node)
		if fci.labelIndex >= 0 {
			label = node.Field(fci.labelIndex).String()
		}
		if !labelOnly {
			if fci.propertiesIndex >= 0 {
				properties = node.Field(fci.propertiesIndex).Interface().([]string)
			}
			if fci.childrenIndex >= 0 {
				children = node.Field(fci.childrenIndex)
				if v.SortNodes {
					children = sortNodes(v, children, true)
				}
			}
		}
		return
	// Gets the (well-known) key-values for label, properties, and children in
	// a map. Again, all these keys-values are optional and will default to
	// zero if missing.
	case reflect.Map:
		if lbl := node.MapIndex(reflect.ValueOf("label")); lbl.Kind() != reflect.Invalid {
			label = lbl.Interface().(string)
		}
		if pps := node.MapIndex(reflect.ValueOf("properties")); pps.Kind() != reflect.Invalid {
			properties = pps.Elem().Interface().([]string)
		}
		if !labelOnly {
			if chs := node.MapIndex(reflect.ValueOf("children")); chs.Kind() != reflect.Invalid {
				children = chs.Elem()
				if v.SortNodes {
					children = sortNodes(v, children, true)
				}
			} else {
				children = reflect.ValueOf([]interface{}{})
			}
		}
		return
	default:
		panic(fmt.Sprintf("unsupported asciitree node or root type %q: %v", node.Kind(), node))
	}
}

// Data structure to sort a list of children (referenced by its dedicated
// slice) by their labels. In order to speed up repeated label lookups during
// sorting, we store the list of labels. However, we don't store the child
// elements, but instead just reference the slice where the childs are finally
// stored/referenced.
type labelledNodes struct {
	labels []string      // the discovered labels for the nodes (see next).
	nodes  reflect.Value // reference to the nodes slice to be sorted.
}

// Returns number of children (interface sort.Interface).
func (l labelledNodes) Len() int { return len(l.labels) }

// Compares two child nodes by their labels (interface sort.Interface).
func (l labelledNodes) Less(i, j int) bool { return l.labels[i] < l.labels[j] }

// Swaps two child nodes (interface sort.Interface). Swapping the child references
// that are in the form of reflect.Values is unfortunately slightly involved, as
// a simple swap without an intermediate temporary would fail.
func (l labelledNodes) Swap(i, j int) {
	l.labels[i], l.labels[j] = l.labels[j], l.labels[i]
	// We cannot simply swap two elements in a slice if they are more
	// intricate types, such as strings, interfaces, maps, et cetera, as
	// opposed to ints. See also: https://github.com/golang/go/issues/3126
	// Instead, we need to dance around and sacrifice the Go(ds) of
	// reflection.
	temp := reflect.New(l.nodes.Index(i).Type()).Elem()
	temp.Set(l.nodes.Index(i))
	l.nodes.Index(i).Set(l.nodes.Index(j))
	l.nodes.Index(j).Set(temp)
}

// Sorts the children by their labels. Optionally works on a copy, if needed.
func sortNodes(v *MapStructVisitor, nodes reflect.Value, copy bool) reflect.Value {
	nodes = reflect.Indirect(nodes)
	// We don't need to dance around like mad when there's nothing to sort. In
	// this case, leave the party early.
	count := nodes.Len()
	if count == 0 {
		return nodes
	}
	// As we are going to sort in place, in some situations we will be asked
	// by the caller to work on a copy instead, because otherwise we would
	// mess with the user's data.
	if copy {
		n := reflect.MakeSlice(nodes.Type(), count, count)
		reflect.Copy(n, nodes)
		nodes = n
	}
	// Fetch all labels for the nodes slice and then sort the nodes slice
	// (this will happen "in place" with respect to the nodes slice).
	list := labelledNodes{labels: make([]string, count), nodes: nodes}
	for idx := 0; idx < count; idx++ {
		elem, ok := nodes.Index(idx).Interface().(reflect.Value)
		if ok {
			list.labels[idx] = v.Label(elem)
		} else {
			list.labels[idx] = v.Label(nodes.Index(idx))
		}
	}
	sort.Sort(list)
	return nodes
}
