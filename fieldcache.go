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
	"slices"
	"sync"
)

// structFields caches certain field indices relevant to rendering trees from
// structs, such as the label field, the properties and children fields, and (if
// any) the roots field.
type structFields struct {
	LabelPath      []int // indices path of label field, or nil.
	PropertiesPath []int // indices path of properties field, or nil.
	ChildrenPath   []int // indices path of children field, or nil.
	RootsPath      []int // indices path of roots field, or nil.
}

// structFieldsCache is our program-global cache for quickly looking up the
// relevant field indices for a particular type.
var structFieldsCache sync.Map

// Returns the field indices for tagged structs, based on a specific node type.
// We employ caching in order to avoid finding the fields (field indices) over
// and over again, especially for mono-type struct trees.
func structFieldInfo(node reflect.Value) *structFields {
	return structInfoCache(&structFieldsCache, node)
}

func structInfoCache(cache *sync.Map, node reflect.Value) *structFields {
	if node.Kind() != reflect.Struct {
		return nil
	}
	// Try to look up this (struct) type from the cache, if already known.
	structT := node.Type()
	if sf, ok := cache.Load(structT); ok {
		return sf.(*structFields)
	}
	// This struct type is yet unknown, so scan the type's fields for
	// asciitree tags, and if found and valid, then learn the field indices.
	newsf := &structFields{}
	findFieldsRecursively(structT, nil, newsf)
	sf, _ := cache.LoadOrStore(structT, newsf)
	return sf.(*structFields)
}

// findsFieldsRecursively locates fields marked as label, properties, children,
// and roots fields, recording their indices paths in the referenced
// structFields value. It recursively descends into anonymous structures fields,
// in a depth first manner, but it does not descend into any named structure
// fields.
func findFieldsRecursively(structT reflect.Type, path []int, sf *structFields) {
	for fieldIdx := range structT.NumField() {
		field := structT.Field(fieldIdx)
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			// we need to dig deeper; please note that this is a non-concurrent
			// use of the path parameter, so we're safe to just use append here
			// without explicit cloning first, as it is fine to reuse the
			// backing array.
			findFieldsRecursively(field.Type, append(path, fieldIdx), sf)
			continue
		}
		if sf.LabelPath == nil && hasAsciitreeTagValue(field, "label") {
			sf.LabelPath = append(slices.Clone(path), fieldIdx)
			continue
		}
		if sf.PropertiesPath == nil && hasAsciitreeTagValue(field, "properties") {
			sf.PropertiesPath = append(slices.Clone(path), fieldIdx)
			continue
		}
		if sf.ChildrenPath == nil && hasAsciitreeTagValue(field, "children") {
			sf.ChildrenPath = append(slices.Clone(path), fieldIdx)
			continue
		}
		if sf.RootsPath == nil && hasAsciitreeTagValue(field, "roots") {
			sf.RootsPath = append(slices.Clone(path), fieldIdx)
			continue
		}
	}
}

// hasAsciitreeTagValue returns true, if the passed field has the "asciitree" tag
// with the specified value; otherwise false.
func hasAsciitreeTagValue(field reflect.StructField, value string) bool {
	v, ok := field.Tag.Lookup("asciitree")
	return ok && v == value
}
