package btree

import "errors"

// ErrNodeNotFound is returned when the cacher could not find a node.
var ErrNodeNotFound = errors.New(`node not found`)

// ErrTreeNotFound is returned when a tree with the provided key could
// not be loaded.
var ErrTreeNotFound = errors.New(`tree not found`)
