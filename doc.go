/*

Package asciitree pretty-prints hierarchical node data structures as ASCII
trees. Asciitree offers different styles, "pure ASCII" and Unicode-based line
"graphics" (see illustration below).


    root1
    ├── 1
    ├── 2
    │   ├── 2.1
    │   └── 2.2
    └── 3
        └── 3.1
    root2
    └── X


User-defined tree data structures can automatically be traversed if they are
either structs or maps and they have been tagged. In many situations this
avoids having to write adaptors (or more specific: visitors) that adapt your
user data to the form needed by asciitree. With simple tagging, all you need
to do is tag your data structure, such as in:

    type node struct {
        Label    string   `asciitree:"label"`
        Props    []string `asciitree:"properties"`
        Children []node   `asciitree:"children"`
    }

Asciitree can both work with a single-root tree, as well as with multiple
roots. Simply pass the Render() function either a single root node, or a slice
of root nodes, and it will handle both cases automatically. You can also pass
a user-data struct with only an `asciitree:"roots"` tag, attached to a struct
field storing your root nodes. Traversal will then proceed as usual.

In addition to automatically rendering tagged user-data structs, asciitree can
also automatically traverse maps if those follow these rules: (1) your map
must be of type "map[string]interface{}". And (2), your map needs to use the
well-known map keys "label", "properties", and "children". If one or more of
these keys is missing, asciitree will assume them to be zero. Finally, you can
optionally pass in a top-level map with the well-known map key "roots" holding
your root nodes. Or you can pass in a slide of root nodes. The Render()
function will detect these use case automatically and handle them accordingly.

*/
package asciitree
