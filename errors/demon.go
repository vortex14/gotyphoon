package errors

import Errors "errors"

var DemonFoolish = Errors.New("demon is foolish")
var DemonNotFound = Errors.New("demon not found")
var DemonHasNotProject = Errors.New("demon has not project")
var DemonExecutingWithoutSettings = Errors.New("demon without settings")
