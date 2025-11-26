# asciitree

[![Go Reference](https://pkg.go.dev/badge/github.com/thediveo/go-asciitree.svg)](https://pkg.go.dev/github.com/thediveo/go-asciitree/v2)
[![License](https://img.shields.io/github/license/thediveo/go-asciitree)](https://img.shields.io/github/license/thediveo/go-asciitree)
![build and test](https://github.com/thediveo/go-asciitree/actions/workflows/buildandtest.yaml/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/thediveo/go-asciitree/v2)](https://goreportcard.com/report/github.com/thediveo/go-asciitree/v2)
![Coverage](https://img.shields.io/badge/Coverage-98.2%25-brightgreen)

`asciitree/v2` is a Go package for simple pretty-printing of tree-like
data structures using pure ASCII "edges" or alternatively Unicode characters
for drawing branches and edges.

    root1
    ├── 1
    ├── 2
    │   ├── 2.1
    │   └── 2.2
    └── 3
        └── 3.1
    root2
    └── X

Nodes can optionally be sorted by their labels. In addition, nodes may have
properties (these are flat, so no properties of properties). These properties
can also optionally be sorted.

## Changes in v2

With v1 dating back to 2019 there surely was merit to align v2 better with
today's "idiomatic Go".
- the `Visitor` interface now strictly uses `any` and `[]any` as parameters and
  results types. This should make refactoring existing and writing new visitors
  more straighforward without the cognitive load of juggling `reflect.Value` and
  friends. 
- as part of the API and code renovation, `interface{}` was replaced with `any`.

## DevContainer

> [!CAUTION]
>
> Do **not** use VSCode's "~~Dev Containers: Clone Repository in Container
> Volume~~" command, as it is utterly broken by design, ignoring
> `.devcontainer/devcontainer.json`.

1. `git clone https://github.com/thediveo/go-asciitree`
2. in VSCode: Ctrl+Shift+P, "Dev Containers: Open Workspace in Container..."
3. select `asciitree.code-workspace` and off you go...

## Supported Go Versions

`clippy` supports versions of Go that are noted by the [Go release
policy](https://golang.org/doc/devel/release.html#policy), that is, major
versions _N_ and _N_-1 (where _N_ is the current major version).

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md).

## Copyright and License

`asciitree` is Copyright 2018‒2025 Harald Albrecht, and licensed under the
Apache License, Version 2.0.
