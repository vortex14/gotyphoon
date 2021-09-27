package errors

import Errors "errors"

var MongoNotFoundDB = Errors.New("not found mongo db")
var MongoNotFoundDBMap = Errors.New("not found mongo db map. need call .init")
