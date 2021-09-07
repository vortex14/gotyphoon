package errors

import "errors"

var GraphNameNotFound         = errors.New("graph name not found")
var GraphMainGraphNotFound    = errors.New("main graph not found")
var GraphOptionsNotFound      = errors.New("graph options not found")
var GraphNodeOptionsNotFound  = errors.New("graph node options not found")
var GraphOptionsLabelRequired = errors.New("graph options label required")
var GraphParentNodeNotFound   = errors.New("graph parent for current node not found ")
var GraphEdgeContextBroken    = errors.New("edge context is broken. required node1 -> node2 and edgeOptions ")
