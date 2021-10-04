package errors

import Errors "errors"

var JqExecuteQueryError = Errors.New("jq query has error. need check input data and jq rule")