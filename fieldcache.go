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

// Caches the label and children field indices for a specific reflection type.
// Oh, and properties, and ... roots (to unify things slightly).
type fieldCacheItem struct {
	labelIndex      int
	propertiesIndex int
	childrenIndex   int
	rootsIndex      int
}

// Cache for quickly looking up the label, children field indices for a given
// reflection type.
var fieldCache = make(map[reflect.Type]*fieldCacheItem)

// Returns the field indices for tagged structs, based on a specific node. We
// employ caching in order to avoid finding the fields (field indices) over
// and over again, especially for mono-type struct trees.
func structInfo(node reflect.Value) (sinfo *fieldCacheItem) {
	if node.Kind() != reflect.Struct {
		return nil
	}
	// Try to look up this (struct) type from the cache, if already known.
	typ := node.Type()
	sinfo, found := fieldCache[typ]
	if found {
		return
	}
	// This struct type is yet unknown, so scan the type's fields for
	// asciitree tags, and if found and valid, then learn the field indices.
	sinfo = &fieldCacheItem{
		labelIndex:      -1,
		propertiesIndex: -1,
		childrenIndex:   -1,
		rootsIndex:      -1,
	}
	for idx := 0; idx < typ.NumField(); idx++ {
		field := typ.Field(idx)
		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			// Oops, it's an anonymous embedded struct, so we need to check that too!
			anon := structInfo(node.Field(idx))
			if anon.labelIndex >= 0 {
				if sinfo.labelIndex >= 0 {
					panic(fmt.Sprintf("double ascii:\"label\" tag for anonymously embedded type %T", field))
				}
				sinfo.labelIndex = anon.labelIndex
			}
			if anon.propertiesIndex >= 0 {
				if sinfo.propertiesIndex >= 0 {
					panic(fmt.Sprintf("double ascii:\"properties\" tag for anonymously embedded type %T", field))
				}
				sinfo.propertiesIndex = anon.propertiesIndex
			}
			if anon.childrenIndex >= 0 {
				if sinfo.childrenIndex >= 0 {
					panic(fmt.Sprintf("double ascii:\"children\" tag for anonymously embedded type %T", field))
				}
				sinfo.childrenIndex = anon.childrenIndex
			}
			if anon.rootsIndex >= 0 {
				if sinfo.rootsIndex >= 0 {
					panic(fmt.Sprintf("double ascii:\"roots\" tag for anonymously embedded type %T", field))
				}
				sinfo.rootsIndex = anon.rootsIndex
			}
		}
		tags, ok := asciitreeTagValues(typ.Field(idx).Tag.Get("asciitree"))
		if !ok {
			panic(fmt.Sprintf("invalid asciitree tag(s) %v", tags))
		}
		for _, tag := range tags {
			switch tag {
			case "label":
				if sinfo.labelIndex >= 0 {
					panic(fmt.Sprintf("double ascii:\"label\" tag for type %T", node))
				}
				sinfo.labelIndex = idx
			case "properties":
				if sinfo.propertiesIndex >= 0 {
					panic(fmt.Sprintf("double ascii:\"properties\" tag for type %T", node))
				}
				sinfo.propertiesIndex = idx
			case "children":
				if sinfo.childrenIndex >= 0 {
					panic(fmt.Sprintf("double ascii:\"children\" tag for type %T", node))
				}
				sinfo.childrenIndex = idx
			case "roots":
				if sinfo.rootsIndex >= 0 {
					panic(fmt.Sprintf("double ascii:\"roots\" tag for type %T", node))
				}
				sinfo.rootsIndex = idx
			}
		}
	}
	fieldCache[typ] = sinfo // cache it.
	return
}

// Returns the (split) tag values of an asciitree tag if they are valid,
// together with an "ok" indication. Otherwise, returns the invalid tag values
// (and only those) with a "nok".
func asciitreeTagValues(tagval string) ([]string, bool) {
	vals := strings.Split(tagval, ",")
	errors := []string{}
	values := make([]string, 0, len(vals))
	for _, value := range vals {
		value = strings.TrimSpace(value)
		if len(value) > 0 {
			switch value {
			case "roots", "label", "properties", "children":
				values = append(values, value)
			default:
				errors = append(errors, value)
			}
		}
	}
	if len(errors) > 0 {
		return errors, false
	}
	return values, true
}
