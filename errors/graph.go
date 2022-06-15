package errors

import "errors"

var GraphNameNotFound         = errors.New("graph name not found")
var GraphResourceNotFound     = errors.New("graph resource not found ")
var GraphMainGraphNotFound    = errors.New("main graph not found")
var GraphNotFoundFormat       = errors.New("not found graph export format")

var GraphEdgeContextBroken    = errors.New("edge context is broken. required node1 -> node2 and edgeOptions ")
var GraphEdgeOptionsNotFound  = errors.New("graph edge options not found ")
var GraphEdgeNotFound         = errors.New("graph edge not found")

var GraphParentNodeNotFound   = errors.New("graph parent for current node not found ")
var GraphNodeOptionsNotFound  = errors.New("graph node options not found")

var GraphOptionsLabelRequired = errors.New("graph options label required")
var GraphOptionsNotFound      = errors.New("graph options not found")

var GraphResourceContextInvalid = errors.New("graph resource context invalid")
var GraphActionContextInvalid   = errors.New("graph action context invalid")
